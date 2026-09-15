package dto

import (
	"testing"

	"github.com/mahditd/zarrine-baft-backend/internal/domain/models"
)

func TestFromProductRequest_PrefersSnapshotOverPreloaded(t *testing.T) {
	req := &models.ProductRequest{
		ID:            1,
		RequestNumber: "2026-000001",
		Status:        models.RequestNew,
		Items: []models.ProductRequestItem{
			{
				ID:               1,
				ProductVariantID: 10,
				Quantity:         2,
				PriceSnapshot:    100,
				// Snapshot (old truth) differs from current preloaded data (edited later)
				ProductCode:   "001",
				ProductNameFA: "نام قدیمی",
				ProductNameEN: "Old Name",
				ColorNameFA:   "آبی",
				ColorNameEN:   "Blue",
				SizeName:      "L",
				ProductVariant: &models.ProductVariant{
					Product: &models.Product{ProductCode: "999", NameFA: "نام جدید", NameEN: "New Name"},
					Color:   &models.Color{NameFA: "قرمز", NameEN: "Red"},
					Size:    &models.Size{Name: "XL"},
				},
			},
		},
	}

	res := FromProductRequestForAdmin(req)
	if len(res.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(res.Items))
	}
	item := res.Items[0]
	if item.ProductCode != "001" || item.ProductNameFA != "نام قدیمی" || item.ProductNameEN != "Old Name" {
		t.Errorf("expected snapshot product preserved, got %+v", item)
	}
	if item.ColorNameFA != "آبی" || item.SizeName != "L" {
		t.Errorf("expected snapshot color/size preserved, got %+v", item)
	}
}

func TestFromProductRequest_FallbackForLegacyItems(t *testing.T) {
	req := &models.ProductRequest{
		ID:            2,
		RequestNumber: "2026-000002",
		Status:        models.RequestNew,
		Items: []models.ProductRequestItem{
			{
				ID:               2,
				ProductVariantID: 11,
				Quantity:         1,
				PriceSnapshot:    200,
				// No snapshot (legacy row) -> fallback to preload
				ProductVariant: &models.ProductVariant{
					Product: &models.Product{ProductCode: "005", NameFA: "کت", NameEN: "Coat"},
					Color:   &models.Color{NameFA: "سبز", NameEN: "Green"},
					Size:    &models.Size{Name: "M"},
				},
			},
		},
	}

	res := FromProductRequestForCustomer(req)
	item := res.Items[0]
	if item.ProductCode != "005" || item.ProductNameFA != "کت" {
		t.Errorf("expected fallback to preloaded product, got %+v", item)
	}
	if item.ColorNameEN != "Green" || item.SizeName != "M" {
		t.Errorf("expected fallback to preloaded color/size, got %+v", item)
	}
}
