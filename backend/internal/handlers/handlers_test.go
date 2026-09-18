package handlers

import (
	"encoding/json"
	"testing"

	"digcatalog/internal/models"
)

func TestFindReqMaterialIDPresence(t *testing.T) {
	t.Run("absent means keep", func(t *testing.T) {
		var req findReq
		if err := json.Unmarshal([]byte(`{"unitId":1,"registerNo":"A-1"}`), &req); err != nil {
			t.Fatal(err)
		}
		if req.MaterialIDSet {
			t.Fatal("MaterialIDSet should be false when materialId is absent")
		}
	})

	t.Run("explicit null means clear", func(t *testing.T) {
		var req findReq
		if err := json.Unmarshal([]byte(`{"materialId":null}`), &req); err != nil {
			t.Fatal(err)
		}
		if !req.MaterialIDSet {
			t.Fatal("MaterialIDSet should be true when materialId is explicitly null")
		}
		if req.MaterialID != nil {
			t.Fatal("MaterialID should be nil for explicit null")
		}
	})

	t.Run("value means set", func(t *testing.T) {
		var req findReq
		if err := json.Unmarshal([]byte(`{"materialId":3}`), &req); err != nil {
			t.Fatal(err)
		}
		if !req.MaterialIDSet {
			t.Fatal("MaterialIDSet should be true when materialId is provided")
		}
		if req.MaterialID == nil || *req.MaterialID != 3 {
			t.Fatalf("MaterialID should be 3, got %v", req.MaterialID)
		}
	})
}

func TestApplyFindReqKeepsMaterialWhenAbsent(t *testing.T) {
	origID := uint(7)
	find := models.Find{MaterialID: &origID, MaterialName: "陶器"}
	req := findReq{UnitID: 1, RegisterNo: "A-1", ArtifactType: "陶片"} // MaterialIDSet=false

	h := &Handler{} // 不触达 DB：未携带 materialId 时应直接返回
	h.applyFindReq(&find, &req)

	if find.MaterialID == nil || *find.MaterialID != origID {
		t.Fatalf("MaterialID should be kept as %d, got %v", origID, find.MaterialID)
	}
	if find.MaterialName != "陶器" {
		t.Fatalf("MaterialName should be kept, got %q", find.MaterialName)
	}
}

func TestApplyFindReqClearsMaterialOnExplicitNull(t *testing.T) {
	origID := uint(7)
	find := models.Find{MaterialID: &origID, MaterialName: "陶器"}
	req := findReq{UnitID: 1, MaterialID: nil, MaterialIDSet: true}

	h := &Handler{}
	h.applyFindReq(&find, &req)

	if find.MaterialID != nil {
		t.Fatalf("MaterialID should be cleared, got %v", find.MaterialID)
	}
	if find.MaterialName != "" {
		t.Fatalf("MaterialName should be cleared, got %q", find.MaterialName)
	}
}

func TestApplyFindReqUpdatesOtherFields(t *testing.T) {
	find := models.Find{Description: "旧描述", StorageLoc: "旧库房"}
	req := findReq{
		UnitID: 2, RegisterNo: "B-2", ArtifactType: "玉器",
		Completeness: "残缺", Description: "新描述", StorageLoc: "新库房",
	}

	h := &Handler{}
	h.applyFindReq(&find, &req)

	if find.Description != "新描述" || find.StorageLoc != "新库房" {
		t.Fatalf("other fields should update: %+v", find)
	}
	if find.UnitID != 2 || find.RegisterNo != "B-2" || find.ArtifactType != "玉器" || find.Completeness != "残缺" {
		t.Fatalf("other fields should update: %+v", find)
	}
}
