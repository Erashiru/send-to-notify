package models

import "time"

type Category struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Description  any    `json:"description"`
	Order        int    `json:"order"`
	IsActive     bool   `json:"isActive"`
	SubCatOrders []int  `json:"subCatOrders"`
	Aliases      []struct {
		Value string `json:"value"`
		Name  string `json:"name"`
		Code  string `json:"code"`
		Dir   string `json:"dir"`
	} `json:"aliases"`
	Extra struct {
		AliasName              string `json:"aliasName"`
		AliasDescription       any    `json:"aliasDescription"`
		DescriptionTranslation any    `json:"descriptionTranslation"`
		IsEgg                  bool   `json:"isEgg"`
	} `json:"extra"`
	SubCategories []struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description any    `json:"description"`
		Order       int    `json:"order"`
		IsActive    bool   `json:"isActive"`
		IsAddons    bool   `json:"isAddons"`
		Aliases     []struct {
			Value string `json:"value"`
			Name  string `json:"name"`
			Code  string `json:"code"`
			Dir   string `json:"dir"`
		} `json:"aliases"`
		Extra struct {
			AliasName              string `json:"aliasName"`
			AliasDescription       any    `json:"aliasDescription"`
			DescriptionTranslation any    `json:"descriptionTranslation"`
			IsEgg                  bool   `json:"isEgg"`
		} `json:"extra"`
		NotInBilling bool  `json:"notInBilling"`
		EntityOrders []int `json:"entityOrders"`
		Entities     []struct {
			NutritionalInfo struct {
				ProteinCount struct {
					Value any    `json:"value"`
					Unit  string `json:"unit"`
				} `json:"proteinCount"`
				FatCount struct {
					Value any    `json:"value"`
					Unit  string `json:"unit"`
				} `json:"fatCount"`
				FiberCount struct {
					Value any    `json:"value"`
					Unit  string `json:"unit"`
				} `json:"fiberCount"`
				CarbohydrateCount struct {
					Value any    `json:"value"`
					Unit  string `json:"unit"`
				} `json:"carbohydrateCount"`
				CalorieCount struct {
					Value any    `json:"value"`
					Unit  string `json:"unit"`
				} `json:"calorieCount"`
			} `json:"nutritionalInfo"`
			ID                           string `json:"id"`
			Name                         string `json:"name"`
			SubCategoryID                any    `json:"subCategoryId"`
			SubCategoryNotInBilling      any    `json:"subCategoryNotInBilling"`
			SubCategoryIsAddons          any    `json:"subCategoryIsAddons"`
			SubCategoryAddOnASNormalItem any    `json:"subCategoryAddOnASNormalItem"`
			EntityType                   string `json:"entityType"`
			HasVariant                   bool   `json:"hasVariant"`
			IsVariant                    bool   `json:"isVariant"`
			Order                        int    `json:"order"`
			IsActive                     bool   `json:"isActive"`
			ClusterStatus                bool   `json:"clusterStatus"`
			Duplicated                   bool   `json:"duplicated"`
			Description                  string `json:"description"`
			Extra                        struct {
				AliasName              string `json:"aliasName"`
				AliasDescription       string `json:"aliasDescription"`
				DescriptionTranslation string `json:"descriptionTranslation"`
				IsEgg                  bool   `json:"isEgg"`
			} `json:"_extra"`
			Tags    []string `json:"tags"`
			Aliases []struct {
				Value string `json:"value"`
				Name  string `json:"name"`
				Code  string `json:"code"`
				Dir   string `json:"dir"`
			} `json:"aliases"`
			Price    float64 `json:"price"`
			IsVeg    bool    `json:"isVeg"`
			TableTab []struct {
				TabID      string  `json:"tabId"`
				ItemRate   float64 `json:"itemRate"`
				ItemStatus string  `json:"itemStatus"`
			} `json:"tableTab"`
			InStock       any      `json:"inStock"`
			PackingCharge any      `json:"packingCharge"`
			Modifiers     []string `json:"modifiers"`
			Variants      []any    `json:"variants"`
			Taxes         []struct {
				ID             string `json:"_id"`
				Name           string `json:"name"`
				TaxCode        any    `json:"taxCode"`
				Type           string `json:"type"`
				Value          int    `json:"value"`
				IsCharge       bool   `json:"isCharge"`
				IsNonRemovable bool   `json:"isNonRemovable"`
				IsGSTOSC       bool   `json:"isGSTOSC"`
				IsGST          bool   `json:"isGST"`
				ApplicableOn   any    `json:"applicableOn"`
				CascadingTaxes []any  `json:"cascadingTaxes"`
				IsAggregator   bool   `json:"isAggregator"`
			} `json:"taxes"`
			ImageURL        string `json:"image_url"`
			ImageURLPng     string `json:"image_url_png"`
			AggregatorImage []struct {
				AggreagtorName string    `json:"aggreagtorName"`
				Jpg            string    `json:"jpg"`
				Png            string    `json:"png"`
				ReferenceID    string    `json:"referenceId"`
				Updated        time.Time `json:"updated"`
				MultipleImages []any     `json:"multipleImages"`
			} `json:"aggregator_image"`
			IsVariantItem bool   `json:"isVariantItem"`
			Scheduler     any    `json:"scheduler"`
			Allergens     []any  `json:"allergens"`
			ServingInfo   string `json:"servingInfo"`
			ServingSize   struct {
				Value any    `json:"value"`
				Unit  string `json:"unit"`
			} `json:"servingSize"`
			SwiggyNutritionInfo struct {
				SpicyLevel       any   `json:"spicyLevel"`
				SweetLevel       any   `json:"sweetLevel"`
				BoneProperty     any   `json:"boneProperty"`
				GravyProperty    any   `json:"gravyProperty"`
				SeasionIngredent any   `json:"seasionIngredent"`
				ServingPeople    any   `json:"servingPeople"`
				ServingSize      any   `json:"servingSize"`
				Accompaniments   []any `json:"accompaniments"`
			} `json:"swiggyNutritionInfo"`
			CalorieCount any  `json:"calorieCount"`
			IsDefault    bool `json:"isDefault"`
			ShelfLife    struct {
				Value any    `json:"value"`
				Unit  string `json:"unit"`
			} `json:"shelfLife"`
			PiecesPerKg struct {
				Min any `json:"min"`
				Max any `json:"max"`
			} `json:"piecesPerKg"`
			OriginalName string `json:"originalName"`
		} `json:"entities"`
		Scheduler    any    `json:"scheduler"`
		OriginalName string `json:"originalName"`
	} `json:"subCategories"`
	Scheduler    any    `json:"scheduler"`
	OriginalName string `json:"originalName"`
}

type Modifier struct {
	Type         string `json:"type"`
	ID           string `json:"_id"`
	Name         string `json:"name"`
	OriginalName string `json:"original_name"`
	Min          int    `json:"min"`
	Max          int    `json:"max"`
	Order        int    `json:"order"`
	IsActive     bool   `json:"isActive"`
	Aliases      []struct {
		Value string `json:"value"`
		Name  any    `json:"name"`
		Code  string `json:"code"`
		Dir   any    `json:"dir"`
	} `json:"aliases"`
	ConstituentItems []struct {
		NutritionalInfo struct {
			ProteinCount struct {
				Value any    `json:"value"`
				Unit  string `json:"unit"`
			} `json:"proteinCount"`
			FatCount struct {
				Value any    `json:"value"`
				Unit  string `json:"unit"`
			} `json:"fatCount"`
			FiberCount struct {
				Value any    `json:"value"`
				Unit  string `json:"unit"`
			} `json:"fiberCount"`
			CarbohydrateCount struct {
				Value any    `json:"value"`
				Unit  string `json:"unit"`
			} `json:"carbohydrateCount"`
			CalorieCount struct {
				Value any    `json:"value"`
				Unit  string `json:"unit"`
			} `json:"calorieCount"`
		} `json:"nutritionalInfo"`
		ID                           string `json:"id"`
		Name                         string `json:"name"`
		SubCategoryID                string `json:"subCategoryId"`
		SubCategoryNotInBilling      bool   `json:"subCategoryNotInBilling"`
		SubCategoryIsAddons          bool   `json:"subCategoryIsAddons"`
		SubCategoryAddOnASNormalItem bool   `json:"subCategoryAddOnASNormalItem"`
		EntityType                   string `json:"entityType"`
		HasVariant                   bool   `json:"hasVariant"`
		IsVariant                    bool   `json:"isVariant"`
		Order                        int    `json:"order"`
		IsActive                     bool   `json:"isActive"`
		ClusterStatus                bool   `json:"clusterStatus"`
		Duplicated                   bool   `json:"duplicated"`
		Description                  string `json:"description"`
		Extra                        struct {
			AliasName              any  `json:"aliasName"`
			AliasDescription       any  `json:"aliasDescription"`
			DescriptionTranslation any  `json:"descriptionTranslation"`
			IsEgg                  bool `json:"isEgg"`
		} `json:"_extra"`
		Tags    []string `json:"tags"`
		Aliases []struct {
			Value string `json:"value"`
			Name  string `json:"name"`
			Code  string `json:"code"`
			Dir   string `json:"dir"`
		} `json:"aliases"`
		Price    float64 `json:"price"`
		IsVeg    bool    `json:"isVeg"`
		TableTab []struct {
			TabID      string  `json:"tabId"`
			ItemRate   float64 `json:"itemRate"`
			ItemStatus string  `json:"itemStatus"`
		} `json:"tableTab"`
		InStock       bool  `json:"inStock"`
		PackingCharge int   `json:"packingCharge"`
		Modifiers     []any `json:"modifiers"`
		Variants      any   `json:"variants"`
		Taxes         []struct {
			ID             string `json:"_id"`
			Name           string `json:"name"`
			TaxCode        any    `json:"taxCode"`
			Type           string `json:"type"`
			Value          int    `json:"value"`
			IsCharge       bool   `json:"isCharge"`
			IsNonRemovable bool   `json:"isNonRemovable"`
			IsGSTOSC       bool   `json:"isGSTOSC"`
			IsGST          bool   `json:"isGST"`
			ApplicableOn   any    `json:"applicableOn"`
			CascadingTaxes []any  `json:"cascadingTaxes"`
			IsAggregator   bool   `json:"isAggregator"`
		} `json:"taxes"`
		ImageURL        string `json:"image_url"`
		ImageURLPng     string `json:"image_url_png"`
		AggregatorImage []struct {
			AggreagtorName string    `json:"aggreagtorName"`
			Jpg            string    `json:"jpg"`
			Png            string    `json:"png"`
			ReferenceID    string    `json:"referenceId"`
			Updated        time.Time `json:"updated"`
			MultipleImages []any     `json:"multipleImages"`
		} `json:"aggregator_image"`
		IsVariantItem bool  `json:"isVariantItem"`
		Scheduler     any   `json:"scheduler"`
		Allergens     []any `json:"allergens"`
		ServingInfo   any   `json:"servingInfo"`
		ServingSize   struct {
			Value any    `json:"value"`
			Unit  string `json:"unit"`
		} `json:"servingSize"`
		SwiggyNutritionInfo struct {
			SpicyLevel       any `json:"spicyLevel"`
			SweetLevel       any `json:"sweetLevel"`
			BoneProperty     any `json:"boneProperty"`
			GravyProperty    any `json:"gravyProperty"`
			SeasionIngredent any `json:"seasionIngredent"`
			ServingPeople    any `json:"servingPeople"`
			ServingSize      any `json:"servingSize"`
			Accompaniments   any `json:"accompaniments"`
		} `json:"swiggyNutritionInfo"`
		CalorieCount any  `json:"calorieCount"`
		IsDefault    bool `json:"isDefault"`
		ShelfLife    struct {
			Value any    `json:"value"`
			Unit  string `json:"unit"`
		} `json:"shelfLife"`
		PiecesPerKg struct {
			Min any `json:"min"`
			Max any `json:"max"`
		} `json:"piecesPerKg"`
	} `json:"constituentItems"`
	Image                any  `json:"image"`
	MultiplePunchMin     int  `json:"multiplePunchMin"`
	MultiplePunchMax     int  `json:"multiplePunchMax"`
	MultiplePunchMaxItem int  `json:"multiplePunchMaxItem"`
	IsPackagingType      bool `json:"isPackagingType"`
	IsAccordionOpen      bool `json:"isAccordionOpen,omitempty"`
}

type MenuCharge struct {
	ID        string `json:"_id"`
	BrandID   string `json:"brand_id"`
	ClusterID string `json:"cluster_id"`
	TabID     string `json:"tab_id"`
	TenantID  string `json:"tenant_id"`
	Value     struct {
		ChargeID     int64  `json:"chargeId"`
		IsActive     bool   `json:"isActive"`
		ChargeNature string `json:"chargeNature"`
		Type         string `json:"type"`
		SelectedTax  any    `json:"selectedTax"`
		Amount       int    `json:"amount"`
		Deployments  []any  `json:"deployments"`
	} `json:"value"`
	Description string   `json:"description"`
	Name        string   `json:"name"`
	Partners    []string `json:"partners"`
}
type Settings struct {
	DeliveryType  string `json:"delivery_type"`
	MinSubtotal   int    `json:"min_subtotal"`
	PushNutrition struct {
		Zomato bool `json:"zomato"`
		Swiggy bool `json:"swiggy"`
	} `json:"pushNutrition"`
	OverrideMinMax bool `json:"overrideMinMax"`
	SendCategories bool `json:"send_categories"`
	ZomatoMenu     struct {
	} `json:"zomatoMenu"`
	ModifierGroupNamePriority string `json:"modifierGroupNamePriority"`
	CareemBrands              []struct {
		BrandName string `json:"brand_name"`
		BrandID   string `json:"brand_id"`
		State     string `json:"state"`
	} `json:"careemBrands"`
	JahezBrandKey      string `json:"jahezBrandKey"`
	ChefzBrandKey      string `json:"chefzBrandKey"`
	ChefzAccessToken   string `json:"chefzAccessToken"`
	ChefzBrandUsername string `json:"chefzBrandUsername"`
	ChefzBrandPassword string `json:"chefzBrandPassword"`
	AllowMultiplePunch bool   `json:"allow_multiple_punch"`
	ShowAllTabs        bool   `json:"show_all_tabs"`
	CxDetails          struct {
		QrFrontImageURL string `json:"qrFrontImageUrl"`
		QrText1         string `json:"qrText_1"`
		QrText2         string `json:"qrText_2"`
		QrStories       []struct {
			URL   string `json:"url"`
			Items []any  `json:"items"`
		} `json:"qrStories"`
	} `json:"cx_details"`
}

type Menu struct {
	Categories      []Category   `json:"categories"`
	Modifiers       []Modifier   `json:"modifiers"`
	Charges         []MenuCharge `json:"charges"`
	BoxComboMapping []any        `json:"boxComboMapping"`
	Settings        Settings     `json:"settings"`
}
