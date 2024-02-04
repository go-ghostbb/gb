package crud

import (
	"context"
	"fmt"
	gbsha256 "ghostbb.io/gb/crypto/gb_sha256"
	gbstr "ghostbb.io/gb/text/gb_str"
	gbconv "ghostbb.io/gb/util/gb_conv"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
)

func (h *Handler) parseLevel(ctx context.Context) cacheLevel {
	value := ctx.Value(CacheCtxLevelKey)
	if value == nil {
		return CacheNone
	}
	if v, ok := value.(cacheLevel); ok {
		return v
	}
	return CacheNone
}

func (h *Handler) genSFKey(tableName string, sql string, vars ...interface{}) string {
	var buf = strings.Builder{}
	buf.WriteString(sql)
	for _, v := range vars {
		buf.WriteString(fmt.Sprintf(":%s", gbconv.String(v)))
	}
	return gbsha256.Encrypt256(fmt.Sprintf("%s:%s", tableName, buf.String()))
}

func (h *Handler) getPrimaryKeys(db *gorm.DB) []string {
	var (
		pks     = make([]string, 0)
		clauses clause.Clause
		where   clause.Where
		ok      bool
		pkName  string
	)
	if db.Statement.Schema == nil {
		return nil
	}
	clauses, ok = db.Statement.Clauses["WHERE"]
	if !ok {
		return nil
	}
	where, ok = clauses.Expression.(clause.Where)
	if !ok {
		return nil
	}

	for _, field := range db.Statement.Schema.Fields {
		if field.PrimaryKey {
			pkName = field.DBName
			break
		}
	}
	if pkName == "" {
		return nil
	}

	for _, expr := range where.Exprs {
		switch e := expr.(type) {
		case clause.Eq:
			if h.getColNameFromColumn(e.Column) == pkName {
				pks = append(pks, gbconv.String(e.Value))
			}
		case clause.IN:
			if h.getColNameFromColumn(e.Column) == pkName {
				pks = append(pks, gbconv.Strings(e.Values)...)
			}
		case clause.Expr:
			var et = h.getExprType(e)
			if et == exprIn || et == exprEq {
				if h.getColNameFromExpr(e, et) == pkName {
					pks = append(pks, h.getPKsFromExpr(e, et)...)
				}
			}
		}
	}

	return h.uniqueStringSlice(pks)
}

func (h *Handler) getColNameFromColumn(col interface{}) string {
	switch v := col.(type) {
	case string:
		return v
	case clause.Column:
		return v.Name
	default:
		return ""
	}
}

func (h *Handler) getExprType(expr clause.Expr) exprType {
	var (
		sql       string
		haveAndOr bool
	)
	// delete spaces
	sql = gbstr.Replace(gbstr.ToLower(expr.SQL), " ", "")

	// see if sql has more than one clause
	haveAndOr = gbstr.Contains(sql, "and") || gbstr.Contains(sql, "or")

	if !haveAndOr && gbstr.Contains(sql, "=") {
		// possibly "id=?" or "id=123"
		fields := gbstr.Split(sql, "=")
		if fields[1] == "?" || gbstr.IsNumeric(fields[1]) {
			return exprEq
		}
	}

	if !haveAndOr && gbstr.Contains(sql, "in") {
		// possibly "idIN(?)"
		fields := strings.Split(sql, "in")
		if len(fields) == 2 {
			// fields[1]
			// possibly "()" or "(123)"
			if len(fields[1]) > 1 && fields[1][0] == '(' && fields[1][len(fields[1])-1] == ')' {
				return exprIn
			}
		}
	}
	return exprOther
}

func (h *Handler) getColNameFromExpr(expr clause.Expr, et exprType) string {
	var sql = gbstr.Replace(gbstr.ToLower(expr.SQL), " ", "")
	if et == exprIn {
		fields := gbstr.Split(sql, "in")
		return fields[0]
	}
	if et == exprEq {
		fields := gbstr.Split(sql, "=")
		return fields[0]
	}
	return ""
}

func (h *Handler) getPKsFromExpr(expr clause.Expr, et exprType) []string {
	var (
		sql = gbstr.Replace(gbstr.ToLower(expr.SQL), " ", "")
		pks = make([]string, 0)
	)

	if et == exprIn {
		var fields = gbstr.Split(sql, "in")
		if len(fields) == 2 && fields[1][0] == '(' && fields[1][len(fields[1])-1] == ')' {
			// (?) or (123, 456, 159, 687) or (1, ?, 2, ?)
			var (
				vars     = gbstr.Split(fields[1][1:len(fields[1])-1], ",")
				varCount = 0
			)
			for _, val := range vars {
				if val == "?" {
					if len(expr.Vars) > varCount {
						pks = append(pks, gbconv.Strings(expr.Vars)...)
						varCount++
					}
					continue
				}
				if gbstr.IsNumeric(val) {
					pks = append(pks, val)
				}
			}
		}
	}

	if et == exprEq {
		var fields = gbstr.Split(sql, "=")
		if len(fields) == 2 {
			if fields[1] == "?" {
				if len(expr.Vars) > 0 {
					pks = append(pks, gbconv.String(expr.Vars))
				}
			} else if gbstr.IsNumeric(fields[1]) {
				pks = append(pks, fields[1])
			}
		}
	}

	return pks
}

func (h *Handler) uniqueStringSlice(strings []string) []string {
	var (
		res   = make([]string, 0)
		check = make(map[string]struct{})
	)
	for _, s := range strings {
		if _, ok := check[s]; !ok {
			res = append(res, s)
			check[s] = struct{}{}
		}
	}
	return res
}

func (h *Handler) hasOtherClauseExceptPrimaryField(db *gorm.DB) bool {
	var (
		clauses clause.Clause
		where   clause.Where
		ok      bool
		pkName  string
	)
	clauses, ok = db.Statement.Clauses["WHERE"]
	if !ok {
		return false
	}
	where, ok = clauses.Expression.(clause.Where)

	for _, field := range db.Statement.Schema.Fields {
		if field.PrimaryKey {
			pkName = field.DBName
			break
		}
	}

	if len(pkName) == 0 {
		return true // return true to skip cache
	}

	for _, expr := range where.Exprs {
		switch e := expr.(type) {
		case clause.Eq:
			if h.getColNameFromColumn(e.Column) != pkName {
				return true
			}
		case clause.IN:
			if h.getColNameFromColumn(e.Column) != pkName {
				return true
			}
		case clause.Expr:
			var et = h.getExprType(e)
			if et == exprIn || et == exprEq {
				if h.getColNameFromExpr(e, et) != pkName {
					return true
				}
			}
		}
	}

	return false
}
