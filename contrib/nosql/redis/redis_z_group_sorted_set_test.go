package redis_test

import (
	gbredis "ghostbb.io/gb/database/gb_redis"
	gbtest "ghostbb.io/gb/test/gb_test"
	gbrand "ghostbb.io/gb/util/gb_rand"
	gbuid "ghostbb.io/gb/util/gb_uid"
	"testing"
)

func Test_GroupSortedSet_ZADD(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		defer redis.FlushDB(ctx)

		var (
			maxn int = 100000000
			k1       = gbuid.S()
			k1m1     = gbuid.S()
			k1m2     = gbuid.S()

			option  gbredis.ZAddOption
			member1 gbredis.ZAddMember
			member2 gbredis.ZAddMember
		)

		member1 = gbredis.ZAddMember{
			Score:  float64(gbrand.Intn(maxn)),
			Member: k1m1,
		}

		_, err := redis.GroupSortedSet().ZAdd(ctx, k1, &option, member1)
		t.AssertNil(err)

		member2 = gbredis.ZAddMember{
			Score:  float64(gbrand.Intn(1000000)),
			Member: k1m2,
		}
		_, err = redis.GroupSortedSet().ZAdd(ctx, k1, &option, member2)
		t.AssertNil(err)

		_, err = redis.GroupSortedSet().ZScore(ctx, k1, k1m1)
		t.AssertNil(err)

		_, err = redis.GroupSortedSet().ZScore(ctx, k1, k1m2)
		t.AssertNil(err)

		var (
			k2   string = gbuid.S()
			k2m1 string = gbuid.S()
			k2m2 string = gbuid.S()
			k2m3 int    = gbrand.Intn(maxn)
		)

		member3 := gbredis.ZAddMember{
			Score:  float64(gbrand.Intn(maxn)),
			Member: k2m1,
		}

		member4 := gbredis.ZAddMember{
			Score:  float64(gbrand.Intn(maxn)),
			Member: k2m2,
		}

		member5 := gbredis.ZAddMember{
			Score:  float64(gbrand.Intn(maxn)),
			Member: k2m3,
		}

		_, err = redis.GroupSortedSet().ZAdd(ctx, k2, &option, member3, member4, member5)
	})

	// with option
	gbtest.C(t, func(t *gbtest.T) {
		defer redis.FlushDB(ctx)

		var (
			maxn int = 100000000
			k1       = gbuid.S()
			k1m1     = gbuid.S()
		)

		member1 := gbredis.ZAddMember{
			Score:  float64(gbrand.Intn(maxn)),
			Member: k1m1,
		}

		option := gbredis.ZAddOption{}
		_, err := redis.GroupSortedSet().ZAdd(ctx, k1, &option, member1)
		t.AssertNil(err)

		// option XX
		optionXX := &gbredis.ZAddOption{
			XX: true,
		}
		memberXX := gbredis.ZAddMember{
			Score:  float64(gbrand.Intn(maxn)),
			Member: k1m1,
		}
		_, err = redis.GroupSortedSet().ZAdd(ctx, k1, optionXX, memberXX)
		t.AssertNil(err)

		scoreXX, err := redis.GroupSortedSet().ZScore(ctx, k1, memberXX.Member)
		t.AssertNil(err)
		t.AssertEQ(scoreXX, memberXX.Score)

		// option NX
		optionNX := &gbredis.ZAddOption{
			NX: true,
		}
		memberNX := gbredis.ZAddMember{
			Score:  float64(gbrand.Intn(maxn)),
			Member: gbuid.S(),
		}
		_, err = redis.GroupSortedSet().ZAdd(ctx, k1, optionNX, memberNX)
		t.AssertNil(err)

		scoreNXOrigin := memberNX.Score
		memberNX.Score = float64(gbrand.Intn(maxn))
		_, err = redis.GroupSortedSet().ZAdd(ctx, k1, optionNX, memberNX)
		t.AssertNil(err)

		score, err := redis.GroupSortedSet().ZScore(ctx, k1, memberNX.Member)
		t.AssertNil(err)
		t.AssertEQ(score, scoreNXOrigin)

		// option LT
		optionLT := &gbredis.ZAddOption{
			LT: true,
		}
		memberLT := gbredis.ZAddMember{
			Score:  float64(gbrand.Intn(maxn)),
			Member: gbuid.S(),
		}
		_, err = redis.GroupSortedSet().ZAdd(ctx, k1, optionLT, memberLT)
		t.AssertNil(err)

		memberLT.Score += 1
		_, err = redis.GroupSortedSet().ZAdd(ctx, k1, optionLT, memberLT)
		t.AssertNil(err)
		scoreLT, err := redis.GroupSortedSet().ZScore(ctx, k1, memberLT.Member)
		t.AssertLT(scoreLT, memberLT.Score)

		memberLT.Score -= 3
		_, err = redis.GroupSortedSet().ZAdd(ctx, k1, optionLT, memberLT)
		t.AssertNil(err)
		scoreLT, err = redis.GroupSortedSet().ZScore(ctx, k1, memberLT.Member)
		t.AssertEQ(scoreLT, memberLT.Score)

		// option GT
		optionGT := &gbredis.ZAddOption{
			GT: true,
		}
		memberGT := gbredis.ZAddMember{
			Score:  float64(gbrand.Intn(maxn)),
			Member: gbuid.S(),
		}
		_, err = redis.GroupSortedSet().ZAdd(ctx, k1, optionGT, memberGT)
		t.AssertNil(err)

		memberLT.Score -= 1
		_, err = redis.GroupSortedSet().ZAdd(ctx, k1, optionGT, memberGT)
		t.AssertNil(err)
		scoreGT, err := redis.GroupSortedSet().ZScore(ctx, k1, memberLT.Member)
		t.AssertGT(scoreGT, memberLT.Score)

		memberLT.Score += 3
		_, err = redis.GroupSortedSet().ZAdd(ctx, k1, optionGT, memberGT)
		t.AssertNil(err)
		scoreGT, err = redis.GroupSortedSet().ZScore(ctx, k1, memberGT.Member)
		t.AssertEQ(scoreGT, memberGT.Score)

		// option CH
		optionCH := &gbredis.ZAddOption{
			CH: true,
		}
		memberCH := gbredis.ZAddMember{
			Score:  float64(gbrand.Intn(maxn)),
			Member: gbuid.S(),
		}
		_, err = redis.GroupSortedSet().ZAdd(ctx, k1, optionCH, memberCH)
		t.AssertNil(err)

		changed, err := redis.GroupSortedSet().ZAdd(ctx, k1, optionCH, memberCH)
		t.AssertNil(err)
		t.AssertEQ(changed.Val(), int64(0))

		memberCH.Score += 1
		changed, err = redis.GroupSortedSet().ZAdd(ctx, k1, optionCH, memberCH)
		t.AssertNil(err)
		t.AssertEQ(changed.Val(), int64(1))

		memberCH.Member = gbuid.S()
		changed, err = redis.GroupSortedSet().ZAdd(ctx, k1, optionCH, memberCH)
		t.AssertNil(err)
		t.AssertEQ(changed.Val(), int64(1))

		// option INCR
		optionINCR := &gbredis.ZAddOption{
			INCR: true,
		}
		memberINCR := gbredis.ZAddMember{
			Score:  float64(gbrand.Intn(maxn)),
			Member: gbuid.S(),
		}
		_, err = redis.GroupSortedSet().ZAdd(ctx, k1, optionINCR, memberINCR)
		t.AssertNil(err)

		scoreIncrOrigin := memberINCR.Score
		memberINCR.Score += 1
		_, err = redis.GroupSortedSet().ZAdd(ctx, k1, optionINCR, memberINCR)
		t.AssertNil(err)

		scoreINCR, err := redis.GroupSortedSet().ZScore(ctx, k1, memberINCR.Member)
		t.AssertNil(err)
		t.AssertEQ(scoreINCR, memberINCR.Score+scoreIncrOrigin)
	})
}

func Test_GroupSortedSet_ZScore(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		defer redis.FlushDB(ctx)

		var (
			k1   string = gbuid.S()
			m1   string = gbuid.S()
			maxn int    = 1000000

			option *gbredis.ZAddOption = new(gbredis.ZAddOption)
		)

		member := gbredis.ZAddMember{
			Member: m1,
			Score:  float64(gbrand.Intn(maxn)),
		}

		_, err := redis.GroupSortedSet().ZAdd(ctx, k1, option, member)
		t.AssertNil(err)

		score, err := redis.GroupSortedSet().ZScore(ctx, k1, m1)
		t.AssertNil(err)
		t.AssertEQ(score, member.Score)

		m2 := gbuid.S()
		score, err = redis.GroupSortedSet().ZScore(ctx, k1, m2)
		t.AssertNil(err)
		t.AssertEQ(score, float64(0))

		k2 := gbuid.S()
		score, err = redis.GroupSortedSet().ZScore(ctx, k2, m2)
		t.AssertNil(err)
		t.AssertEQ(score, float64(0))
	})
}

func Test_GroupSortedSet_ZIncrBy(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		defer redis.FlushDB(ctx)

		k := gbuid.S()
		m := gbuid.S()

		var incr float64 = 6
		_, err := redis.GroupSortedSet().ZIncrBy(ctx, k, incr, m)
		t.AssertNil(err)

		incr2 := float64(3)
		incredScore, err := redis.GroupSortedSet().ZIncrBy(ctx, k, incr2, m)
		t.AssertNil(err)
		t.AssertEQ(incredScore, incr+incr2)
	})
}

func Test_GroupSortedSet_ZCard(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k                          = gbuid.S()
			option *gbredis.ZAddOption = new(gbredis.ZAddOption)
		)

		rand := gbrand.N(10, 20)
		for i := 1; i <= rand; i++ {
			_, err := redis.GroupSortedSet().ZAdd(ctx, k, option, gbredis.ZAddMember{
				Member: i,
				Score:  float64(i),
			})
			t.AssertNil(err)

			cnt, err := redis.GroupSortedSet().ZCard(ctx, k)
			t.AssertNil(err)
			t.AssertEQ(cnt, int64(i))
		}
	})
}

func Test_GroupSortedSet_ZCount(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		defer redis.FlushDB(ctx)

		var (
			k      string              = gbuid.S()
			option *gbredis.ZAddOption = new(gbredis.ZAddOption)
		)

		min, max := "5", "378"
		memSlice := []int{-6, 3, 7, 9, 100, 500, 666}

		for i := 0; i < len(memSlice); i++ {
			redis.GroupSortedSet().ZAdd(ctx, k, option, gbredis.ZAddMember{
				Member: memSlice[i],
				Score:  float64(memSlice[i]),
			})
		}

		cnt, err := redis.GroupSortedSet().ZCount(ctx, k, min, max)
		t.AssertNil(err)
		t.AssertEQ(cnt, int64(3))

		cnt, err = redis.GroupSortedSet().ZCount(ctx, k, "-inf", max)
		t.AssertNil(err)
		t.AssertEQ(cnt, int64(5))

		cnt, err = redis.GroupSortedSet().ZCount(ctx, k, "-inf", "+inf")
		t.AssertNil(err)
		t.AssertEQ(cnt, int64(len(memSlice)))

		cnt, err = redis.GroupSortedSet().ZCount(ctx, k, "(500", "(567")
		t.AssertNil(err)
		t.AssertEQ(cnt, int64(0))

		cnt, err = redis.GroupSortedSet().ZCount(ctx, k, "(500", "+inf")
		t.AssertNil(err)
		t.AssertEQ(cnt, int64(1))

		cnt, err = redis.GroupSortedSet().ZCount(ctx, k, "(3", "(567")
		t.AssertNil(err)
		t.AssertEQ(cnt, int64(4))
	})
}

func Test_GroupSortedSet_ZRange(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		defer redis.FlushDB(ctx)

		var (
			k                          = gbuid.S()
			option *gbredis.ZAddOption = new(gbredis.ZAddOption)
		)

		slice := []string{"one", "two", "three"}
		for i := 0; i < len(slice); i++ {
			redis.ZAdd(ctx, k, option, gbredis.ZAddMember{
				Member: slice[i],
				Score:  float64(i + 1),
			})
		}

		ret, err := redis.GroupSortedSet().ZRange(ctx, k, 0, -1)
		t.AssertNil(err)
		expected := []string{"one", "two", "three"}
		for i := 0; i < len(ret); i++ {
			t.AssertEQ(ret[i].String(), expected[i])
		}

		ret, err = redis.GroupSortedSet().ZRange(ctx, k, 2, 3)
		t.AssertNil(err)
		expected = []string{"three"}
		for i := 0; i < len(ret); i++ {
			t.AssertEQ(ret[i].String(), expected[i])
		}

		// ret, err = redis.GroupSortedSet().ZRange(ctx, k, 0, -1,
		// 	gbredis.ZRangeOption{WithScores: true})
		// t.AssertNil(err)
		// expectedScore := []interface{}{1, "one", 2, "two", 3, "three"}
		// for i := 0; i < len(ret); i++ {
		// 	t.AssertEQ(ret[i].String(), expectedScore[i])
		// }
	})
}

func Test_GroupSortedSet_ZRevRange(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		defer redis.FlushDB(ctx)

		var (
			k                          = gbuid.S()
			option *gbredis.ZAddOption = new(gbredis.ZAddOption)
		)

		slice := []string{"one", "two", "three"}
		for i := 0; i < len(slice); i++ {
			redis.ZAdd(ctx, k, option, gbredis.ZAddMember{
				Member: slice[i],
				Score:  float64(i + 1),
			})
		}

		ret, err := redis.GroupSortedSet().ZRevRange(ctx, k, 0, -1)
		t.AssertNil(err)
		expected := []interface{}{"three", "two", "one"}
		t.AssertEQ(ret.Slice(), expected)

		ret, err = redis.GroupSortedSet().ZRevRange(ctx, k, 0, 1)
		t.AssertNil(err)
		expected = []interface{}{"three", "two"}
		t.AssertEQ(ret.Slice(), expected)
	})
}

func Test_GroupSortedSet_ZRank(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		defer redis.FlushDB(ctx)

		var (
			k                          = gbuid.S()
			option *gbredis.ZAddOption = new(gbredis.ZAddOption)
		)

		slice := []int64{1, 3, 5, 7, 9}
		for i := 0; i < len(slice); i++ {
			redis.ZAdd(ctx, k, option, gbredis.ZAddMember{
				Member: i,
				Score:  float64(slice[i]),
			})
		}

		rank, err := redis.ZRank(ctx, k, 0)
		t.AssertNil(err)
		t.AssertEQ(rank, int64(0))

		rank, err = redis.ZRank(ctx, k, 3)
		t.AssertNil(err)
		t.AssertEQ(rank, int64(3))

		rank, err = redis.ZRank(ctx, k, 6)
		t.AssertNil(err)
		t.AssertEQ(rank, int64(0))
	})
}

func Test_GroupSortedSet_ZRevRank(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		defer redis.FlushDB(ctx)

		var (
			k                          = gbuid.S()
			option *gbredis.ZAddOption = new(gbredis.ZAddOption)
		)

		slice := []int64{1, 3, 5, 7, 9}
		for i := 0; i < len(slice); i++ {
			redis.ZAdd(ctx, k, option, gbredis.ZAddMember{
				Member: i,
				Score:  float64(slice[i]),
			})
		}

		rank, err := redis.ZRevRank(ctx, k, 0)
		t.AssertNil(err)
		t.AssertEQ(rank, int64(4))

		rank, err = redis.ZRevRank(ctx, k, 3)
		t.AssertNil(err)
		t.AssertEQ(rank, int64(1))

		rank, err = redis.ZRevRank(ctx, k, 9)
		t.AssertNil(err)
		t.AssertEQ(rank, int64(0))

		rank, err = redis.ZRevRank(ctx, k, 6)
		t.AssertNil(err)
		t.AssertEQ(rank, int64(0))
	})
}

func Test_GroupSortedSet_ZRem(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		defer redis.FlushDB(ctx)

		var (
			k      = gbuid.S()
			m      = gbuid.S()
			option = new(gbredis.ZAddOption)
		)

		cnt, err := redis.ZRem(ctx, k, m)
		t.AssertNil(err)
		t.AssertEQ(cnt, int64(0))

		member := gbredis.ZAddMember{
			Member: m,
			Score:  123,
		}
		redis.ZAdd(ctx, k, option, member)

		cnt, err = redis.ZRem(ctx, k, m)
		t.AssertNil(err)
		t.AssertEQ(cnt, int64(1))

		member2 := gbredis.ZAddMember{
			Member: gbuid.S(),
			Score:  456,
		}
		_, err = redis.ZAdd(ctx, k, option, member, member2)
		t.AssertNil(err)

		cnt, err = redis.ZRem(ctx, k, m, "non_exists")
		t.AssertNil(err)
		t.AssertEQ(cnt, int64(1))

		_, err = redis.ZAdd(ctx, k, option, member, member2)
		t.AssertNil(err)

		cnt, err = redis.ZRem(ctx, k, m, member2.Member)
		t.AssertNil(err)
		t.AssertEQ(cnt, int64(2))
	})
}

func Test_GroupSortedSet_ZRemRangeByRank(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		defer redis.FlushDB(ctx)

		var (
			k                          = gbuid.S()
			option *gbredis.ZAddOption = new(gbredis.ZAddOption)
		)

		slice := []int64{1, 3, 5, 7, 9}
		for i := 0; i < len(slice); i++ {
			redis.GroupSortedSet().ZAdd(ctx, k, option, gbredis.ZAddMember{
				Member: i,
				Score:  float64(slice[i]),
			})
		}

		rmd, err := redis.GroupSortedSet().ZRemRangeByRank(ctx, k, 0, 2)
		t.AssertNil(err)
		t.AssertEQ(rmd, int64(3))

		score, err := redis.GroupSortedSet().ZScore(ctx, k, 0)
		t.AssertNil(err)
		t.AssertEQ(score, float64(0))

		score, err = redis.GroupSortedSet().ZScore(ctx, k, 1)
		t.AssertNil(err)
		t.AssertEQ(score, float64(0))

		score, err = redis.GroupSortedSet().ZScore(ctx, k, 2)
		t.AssertNil(err)
		t.AssertEQ(score, float64(0))

		score, err = redis.GroupSortedSet().ZScore(ctx, k, 3)
		t.AssertNil(err)
		t.AssertEQ(score, float64(7))
	})

	gbtest.C(t, func(t *gbtest.T) {
		defer redis.FlushDB(ctx)

		var (
			k                          = gbuid.S()
			option *gbredis.ZAddOption = new(gbredis.ZAddOption)
		)

		slice := []int64{1, 3, 5, 7, 9}
		for i := 0; i < len(slice); i++ {
			redis.GroupSortedSet().ZAdd(ctx, k, option, gbredis.ZAddMember{
				Member: i,
				Score:  float64(slice[i]),
			})
		}

		rmd, err := redis.GroupSortedSet().ZRemRangeByRank(ctx, k, -3, -2)
		t.AssertNil(err)
		t.AssertEQ(rmd, int64(2))

		score, err := redis.GroupSortedSet().ZScore(ctx, k, 2)
		t.AssertNil(err)
		t.AssertEQ(score, float64(0))

		score, err = redis.GroupSortedSet().ZScore(ctx, k, 3)
		t.AssertNil(err)
		t.AssertEQ(score, float64(0))

		score, err = redis.GroupSortedSet().ZScore(ctx, k, 4)
		t.AssertNil(err)
		t.AssertEQ(score, float64(9))
	})

	gbtest.C(t, func(t *gbtest.T) {
		defer redis.FlushDB(ctx)

		var (
			k                          = gbuid.S()
			option *gbredis.ZAddOption = new(gbredis.ZAddOption)
		)

		slice := []int64{1, 3, 5, 7, 9}
		for i := 0; i < len(slice); i++ {
			redis.GroupSortedSet().ZAdd(ctx, k, option, gbredis.ZAddMember{
				Member: i,
				Score:  float64(slice[i]),
			})
		}

		rmd, err := redis.GroupSortedSet().ZRemRangeByRank(ctx, k, 3, -1)
		t.AssertNil(err)
		t.AssertEQ(rmd, int64(2))

		score, err := redis.GroupSortedSet().ZScore(ctx, k, 3)
		t.AssertNil(err)
		t.AssertEQ(score, float64(0))

		score, err = redis.GroupSortedSet().ZScore(ctx, k, 4)
		t.AssertNil(err)
		t.AssertEQ(score, float64(0))

		score, err = redis.GroupSortedSet().ZScore(ctx, k, 1)
		t.AssertNil(err)
		t.AssertEQ(score, float64(3))
	})
}
func Test_GroupSortedSet_ZRemRangeByScore(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		defer redis.FlushDB(ctx)

		var (
			k                          = gbuid.S()
			option *gbredis.ZAddOption = new(gbredis.ZAddOption)
		)

		slice := []int64{1, 3, 5, 7, 9}
		for i := 0; i < len(slice); i++ {
			redis.GroupSortedSet().ZAdd(ctx, k, option, gbredis.ZAddMember{
				Member: i,
				Score:  float64(slice[i]),
			})
		}

		rmd, err := redis.GroupSortedSet().ZRemRangeByScore(ctx, k, "(3", "9")
		t.AssertNil(err)
		t.AssertEQ(rmd, int64(3))

		score, err := redis.GroupSortedSet().ZScore(ctx, k, 3)
		t.AssertNil(err)
		t.AssertEQ(score, float64(0))

		score, err = redis.GroupSortedSet().ZScore(ctx, k, 4)
		t.AssertNil(err)
		t.AssertEQ(score, float64(0))

		score, err = redis.GroupSortedSet().ZScore(ctx, k, 5)
		t.AssertNil(err)
		t.AssertEQ(score, float64(0))

		score, err = redis.GroupSortedSet().ZScore(ctx, k, 1)
		t.AssertNil(err)
		t.AssertEQ(score, float64(3))
	})

	gbtest.C(t, func(t *gbtest.T) {
		defer redis.FlushDB(ctx)

		var (
			k                          = gbuid.S()
			option *gbredis.ZAddOption = new(gbredis.ZAddOption)
		)

		slice := []int64{1, 3, 5, 7, 9}
		for i := 0; i < len(slice); i++ {
			redis.GroupSortedSet().ZAdd(ctx, k, option, gbredis.ZAddMember{
				Member: i,
				Score:  float64(slice[i]),
			})
		}

		rmd, err := redis.GroupSortedSet().ZRemRangeByScore(ctx, k, "3", "(9")
		t.AssertNil(err)
		t.AssertEQ(rmd, int64(3))

		score, err := redis.GroupSortedSet().ZScore(ctx, k, 1)
		t.AssertNil(err)
		t.AssertEQ(score, float64(0))

		score, err = redis.GroupSortedSet().ZScore(ctx, k, 2)
		t.AssertNil(err)
		t.AssertEQ(score, float64(0))

		score, err = redis.GroupSortedSet().ZScore(ctx, k, 3)
		t.AssertNil(err)
		t.AssertEQ(score, float64(0))

		score, err = redis.GroupSortedSet().ZScore(ctx, k, 4)
		t.AssertNil(err)
		t.AssertEQ(score, float64(9))
	})

	gbtest.C(t, func(t *gbtest.T) {
		defer redis.FlushDB(ctx)

		var (
			k                          = gbuid.S()
			option *gbredis.ZAddOption = new(gbredis.ZAddOption)
		)

		slice := []int64{1, 3, 5, 7, 9}
		for i := 0; i < len(slice); i++ {
			redis.GroupSortedSet().ZAdd(ctx, k, option, gbredis.ZAddMember{
				Member: i,
				Score:  float64(slice[i]),
			})
		}

		rmd, err := redis.GroupSortedSet().ZRemRangeByScore(ctx, k, "-inf", "9")
		t.AssertNil(err)
		t.AssertEQ(rmd, int64(5))

		cnt, err := redis.GroupSortedSet().ZCard(ctx, k)
		t.AssertNil(err)
		t.AssertEQ(cnt, int64(0))
	})

	gbtest.C(t, func(t *gbtest.T) {
		defer redis.FlushDB(ctx)

		var (
			k                          = gbuid.S()
			option *gbredis.ZAddOption = new(gbredis.ZAddOption)
		)

		slice := []int64{1, 3, 5, 7, 9}
		for i := 0; i < len(slice); i++ {
			redis.GroupSortedSet().ZAdd(ctx, k, option, gbredis.ZAddMember{
				Member: i,
				Score:  float64(slice[i]),
			})
		}

		rmd, err := redis.GroupSortedSet().ZRemRangeByScore(ctx, k, "-inf", "+inf")
		t.AssertNil(err)
		t.AssertEQ(rmd, int64(5))

		cnt, err := redis.GroupSortedSet().ZCard(ctx, k)
		t.AssertNil(err)
		t.AssertEQ(cnt, int64(0))
	})
}
func Test_GroupSortedSet_ZRemRangeByLex(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		defer redis.FlushDB(ctx)

		var (
			k                          = gbuid.S()
			option *gbredis.ZAddOption = new(gbredis.ZAddOption)
		)

		slice := []string{"aaaa", "b", "c", "d", "e", "foo", "zap", "zip", "ALPHA", "alpha"}
		for i := 0; i < len(slice); i++ {
			redis.GroupSortedSet().ZAdd(ctx, k, option, gbredis.ZAddMember{
				Member: slice[i],
				Score:  float64(0),
			})
		}

		cnt, err := redis.GroupSortedSet().ZRemRangeByLex(ctx, k, "[alpha", "[omega")
		t.AssertNil(err)
		t.AssertEQ(cnt, int64(6))

		cnt, err = redis.GroupSortedSet().ZCard(ctx, k)
		t.AssertNil(err)
		t.AssertEQ(cnt, int64(4))
	})
}
func Test_GroupSortedSet_ZLexCount(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		defer redis.FlushDB(ctx)

		var (
			k                          = gbuid.S()
			option *gbredis.ZAddOption = new(gbredis.ZAddOption)
		)

		slice := []string{"a", "b", "c", "d", "e", "f", "g"}
		for i := 0; i < len(slice); i++ {
			redis.GroupSortedSet().ZAdd(ctx, k, option, gbredis.ZAddMember{
				Member: slice[i],
				Score:  float64(0),
			})
		}

		cnt, err := redis.GroupSortedSet().ZLexCount(ctx, k, "-", "+")
		t.AssertNil(err)
		t.AssertEQ(cnt, int64(7))

		cnt, err = redis.GroupSortedSet().ZLexCount(ctx, k, "[b", "[f")
		t.AssertNil(err)
		t.AssertEQ(cnt, int64(5))

	})
}
