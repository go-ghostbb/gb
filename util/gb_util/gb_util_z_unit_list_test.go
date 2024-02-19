package gbutil_test

import (
	"ghostbb.io/gb/frame/g"
	gbtest "ghostbb.io/gb/test/gb_test"
	gbutil "ghostbb.io/gb/util/gb_util"
	"testing"
)

func Test_ListItemValues_Map(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		listMap := g.List{
			g.Map{"id": 1, "score": 100},
			g.Map{"id": 2, "score": 99},
			g.Map{"id": 3, "score": 99},
		}
		t.Assert(gbutil.ListItemValues(listMap, "id"), g.Slice{1, 2, 3})
		t.Assert(gbutil.ListItemValues(&listMap, "id"), g.Slice{1, 2, 3})
		t.Assert(gbutil.ListItemValues(listMap, "score"), g.Slice{100, 99, 99})
	})
	gbtest.C(t, func(t *gbtest.T) {
		listMap := g.List{
			g.Map{"id": 1, "score": 100},
			g.Map{"id": 2, "score": nil},
			g.Map{"id": 3, "score": 0},
		}
		t.Assert(gbutil.ListItemValues(listMap, "id"), g.Slice{1, 2, 3})
		t.Assert(gbutil.ListItemValues(listMap, "score"), g.Slice{100, nil, 0})
	})
	gbtest.C(t, func(t *gbtest.T) {
		listMap := g.List{}
		t.Assert(len(gbutil.ListItemValues(listMap, "id")), 0)
	})
}

func Test_ListItemValues_Map_SubKey(t *testing.T) {
	type Scores struct {
		Math    int
		English int
	}
	gbtest.C(t, func(t *gbtest.T) {
		listMap := g.List{
			g.Map{"id": 1, "scores": Scores{100, 60}},
			g.Map{"id": 2, "scores": Scores{0, 100}},
			g.Map{"id": 3, "scores": Scores{59, 99}},
		}
		t.Assert(gbutil.ListItemValues(listMap, "scores", "Math"), g.Slice{100, 0, 59})
		t.Assert(gbutil.ListItemValues(listMap, "scores", "English"), g.Slice{60, 100, 99})
		t.Assert(gbutil.ListItemValues(listMap, "scores", "PE"), g.Slice{})
	})
}

func Test_ListItemValues_Map_Array_SubKey(t *testing.T) {
	type Scores struct {
		Math    int
		English int
	}
	gbtest.C(t, func(t *gbtest.T) {
		listMap := g.List{
			g.Map{"id": 1, "scores": []Scores{{1, 2}, {3, 4}}},
			g.Map{"id": 2, "scores": []Scores{{5, 6}, {7, 8}}},
			g.Map{"id": 3, "scores": []Scores{{9, 10}, {11, 12}}},
		}
		t.Assert(gbutil.ListItemValues(listMap, "scores", "Math"), g.Slice{1, 3, 5, 7, 9, 11})
		t.Assert(gbutil.ListItemValues(listMap, "scores", "English"), g.Slice{2, 4, 6, 8, 10, 12})
		t.Assert(gbutil.ListItemValues(listMap, "scores", "PE"), g.Slice{})
	})
}

func Test_ListItemValues_Struct(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		type T struct {
			Id    int
			Score float64
		}
		listStruct := g.Slice{
			T{1, 100},
			T{2, 99},
			T{3, 0},
		}
		t.Assert(gbutil.ListItemValues(listStruct, "Id"), g.Slice{1, 2, 3})
		t.Assert(gbutil.ListItemValues(listStruct, "Score"), g.Slice{100, 99, 0})
	})
	// Pointer items.
	gbtest.C(t, func(t *gbtest.T) {
		type T struct {
			Id    int
			Score float64
		}
		listStruct := g.Slice{
			&T{1, 100},
			&T{2, 99},
			&T{3, 0},
		}
		t.Assert(gbutil.ListItemValues(listStruct, "Id"), g.Slice{1, 2, 3})
		t.Assert(gbutil.ListItemValues(listStruct, "Score"), g.Slice{100, 99, 0})
	})
	// Nil element value.
	gbtest.C(t, func(t *gbtest.T) {
		type T struct {
			Id    int
			Score interface{}
		}
		listStruct := g.Slice{
			T{1, 100},
			T{2, nil},
			T{3, 0},
		}
		t.Assert(gbutil.ListItemValues(listStruct, "Id"), g.Slice{1, 2, 3})
		t.Assert(gbutil.ListItemValues(listStruct, "Score"), g.Slice{100, nil, 0})
	})
}

func Test_ListItemValues_Struct_SubKey(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		type Student struct {
			Id    int
			Score float64
		}
		type Class struct {
			Total    int
			Students []Student
		}
		listStruct := g.Slice{
			Class{2, []Student{{1, 1}, {2, 2}}},
			Class{3, []Student{{3, 3}, {4, 4}, {5, 5}}},
			Class{1, []Student{{6, 6}}},
		}
		t.Assert(gbutil.ListItemValues(listStruct, "Total"), g.Slice{2, 3, 1})
		t.Assert(gbutil.ListItemValues(listStruct, "Students"), `[[{"Id":1,"Score":1},{"Id":2,"Score":2}],[{"Id":3,"Score":3},{"Id":4,"Score":4},{"Id":5,"Score":5}],[{"Id":6,"Score":6}]]`)
		t.Assert(gbutil.ListItemValues(listStruct, "Students", "Id"), g.Slice{1, 2, 3, 4, 5, 6})
	})
	gbtest.C(t, func(t *gbtest.T) {
		type Student struct {
			Id    int
			Score float64
		}
		type Class struct {
			Total    int
			Students []*Student
		}
		listStruct := g.Slice{
			&Class{2, []*Student{{1, 1}, {2, 2}}},
			&Class{3, []*Student{{3, 3}, {4, 4}, {5, 5}}},
			&Class{1, []*Student{{6, 6}}},
		}
		t.Assert(gbutil.ListItemValues(listStruct, "Total"), g.Slice{2, 3, 1})
		t.Assert(gbutil.ListItemValues(listStruct, "Students"), `[[{"Id":1,"Score":1},{"Id":2,"Score":2}],[{"Id":3,"Score":3},{"Id":4,"Score":4},{"Id":5,"Score":5}],[{"Id":6,"Score":6}]]`)
		t.Assert(gbutil.ListItemValues(listStruct, "Students", "Id"), g.Slice{1, 2, 3, 4, 5, 6})
	})
}

func Test_ListItemValuesUnique(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		listMap := g.List{
			g.Map{"id": 1, "score": 100},
			g.Map{"id": 2, "score": 100},
			g.Map{"id": 3, "score": 100},
			g.Map{"id": 4, "score": 100},
			g.Map{"id": 5, "score": 100},
		}
		t.Assert(gbutil.ListItemValuesUnique(listMap, "id"), g.Slice{1, 2, 3, 4, 5})
		t.Assert(gbutil.ListItemValuesUnique(listMap, "score"), g.Slice{100})
	})
	gbtest.C(t, func(t *gbtest.T) {
		listMap := g.List{
			g.Map{"id": 1, "score": 100},
			g.Map{"id": 2, "score": 100},
			g.Map{"id": 3, "score": 100},
			g.Map{"id": 4, "score": 100},
			g.Map{"id": 5, "score": 99},
		}
		t.Assert(gbutil.ListItemValuesUnique(listMap, "id"), g.Slice{1, 2, 3, 4, 5})
		t.Assert(gbutil.ListItemValuesUnique(listMap, "score"), g.Slice{100, 99})
	})
	gbtest.C(t, func(t *gbtest.T) {
		listMap := g.List{
			g.Map{"id": 1, "score": 100},
			g.Map{"id": 2, "score": 100},
			g.Map{"id": 3, "score": 0},
			g.Map{"id": 4, "score": 100},
			g.Map{"id": 5, "score": 99},
		}
		t.Assert(gbutil.ListItemValuesUnique(listMap, "id"), g.Slice{1, 2, 3, 4, 5})
		t.Assert(gbutil.ListItemValuesUnique(listMap, "score"), g.Slice{100, 0, 99})
	})
}

func Test_ListItemValuesUnique_Struct_SubKey(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		type Student struct {
			Id    int
			Score float64
		}
		type Class struct {
			Total    int
			Students []Student
		}
		listStruct := g.Slice{
			Class{2, []Student{{1, 1}, {1, 2}}},
			Class{3, []Student{{2, 3}, {2, 4}, {5, 5}}},
			Class{1, []Student{{6, 6}}},
		}
		t.Assert(gbutil.ListItemValuesUnique(listStruct, "Total"), g.Slice{2, 3, 1})
		t.Assert(gbutil.ListItemValuesUnique(listStruct, "Students", "Id"), g.Slice{1, 2, 5, 6})
	})
	gbtest.C(t, func(t *gbtest.T) {
		type Student struct {
			Id    int
			Score float64
		}
		type Class struct {
			Total    int
			Students []*Student
		}
		listStruct := g.Slice{
			&Class{2, []*Student{{1, 1}, {1, 2}}},
			&Class{3, []*Student{{2, 3}, {2, 4}, {5, 5}}},
			&Class{1, []*Student{{6, 6}}},
		}
		t.Assert(gbutil.ListItemValuesUnique(listStruct, "Total"), g.Slice{2, 3, 1})
		t.Assert(gbutil.ListItemValuesUnique(listStruct, "Students", "Id"), g.Slice{1, 2, 5, 6})
	})
}

func Test_ListItemValuesUnique_Map_Array_SubKey(t *testing.T) {
	type Scores struct {
		Math    int
		English int
	}
	gbtest.C(t, func(t *gbtest.T) {
		listMap := g.List{
			g.Map{"id": 1, "scores": []Scores{{1, 2}, {1, 2}}},
			g.Map{"id": 2, "scores": []Scores{{5, 8}, {5, 8}}},
			g.Map{"id": 3, "scores": []Scores{{9, 10}, {11, 12}}},
		}
		t.Assert(gbutil.ListItemValuesUnique(listMap, "scores", "Math"), g.Slice{1, 5, 9, 11})
		t.Assert(gbutil.ListItemValuesUnique(listMap, "scores", "English"), g.Slice{2, 8, 10, 12})
		t.Assert(gbutil.ListItemValuesUnique(listMap, "scores", "PE"), g.Slice{})
	})
}

func Test_ListItemValuesUnique_Binary_ID(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		listMap := g.List{
			g.Map{"id": []byte{1}, "score": 100},
			g.Map{"id": []byte{2}, "score": 100},
			g.Map{"id": []byte{3}, "score": 100},
			g.Map{"id": []byte{4}, "score": 100},
			g.Map{"id": []byte{4}, "score": 100},
		}
		t.Assert(gbutil.ListItemValuesUnique(listMap, "id"), g.Slice{[]byte{1}, []byte{2}, []byte{3}, []byte{4}})
		t.Assert(gbutil.ListItemValuesUnique(listMap, "score"), g.Slice{100})
	})
}
