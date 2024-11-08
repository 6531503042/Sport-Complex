package migration

import (
	"context"
	"fmt"
	"log"
	"main/config"
	"main/modules/facility"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// ensureNoDuplicateFacility ensures that database names are not duplicated
func ensureNoDuplicateFacility(pctx context.Context, client *mongo.Client, originalName string) (*mongo.Database, error) {
	// Create the sanitized name by removing any unwanted suffixes like '_facility_facility'
	normalizedName := strings.Replace(originalName, "_facility_facility", "_facility", -1)

	// Check if the database already exists
	dbNames, err := client.ListDatabaseNames(pctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list database names: %v", err)
	}

	// If the database already exists, return it
	for _, dbName := range dbNames {
		if dbName == normalizedName {
			log.Printf("Database '%s' already exists, skipping creation.", normalizedName)
			return client.Database(normalizedName), nil
		}
	}

	// If it doesn't exist, create the database
	log.Printf("Creating database: %s", normalizedName)
	db := client.Database(normalizedName)
	return db, nil
}

// SetupFacilities creates the required facilities without duplicating the database
func SetupFacilities(pctx context.Context, cfg *config.Config, client *mongo.Client) error {
	facilityNames := []string{"fitness_facility", "swimming_facility", "football_facility", "badminton_facility"}

	for _, facilityName := range facilityNames {
		// Ensure the database name is not duplicated
		db, err := ensureNoDuplicateFacility(pctx, client, facilityName)
		if err != nil {
			return err
		}

		// Create "facility" collection with initial data if it doesn't exist
		if err := createFacilityCollection(pctx, db, facilityName); err != nil {
			log.Fatalf("Failed to create facility collection for %s: %v", facilityName, err)
		}

		// Create slots collection based on facility type
		if facilityName == "badminton_facility" {
			if err := createBadmintonSlots(pctx, db); err != nil {
				log.Fatalf("Failed to create badminton slots for %s: %v", facilityName, err)
			}
			if err := createBadmintonCourts(pctx, db); err != nil {
				log.Fatalf("Failed to create badminton courts for %s: %v", facilityName, err)
			}
		} else {
			if err := createNormalSlots(pctx, db, facilityName); err != nil {
				log.Fatalf("Failed to create normal slots for %s: %v", facilityName, err)
			}
		}
	}

	return nil
}


// Create facility collection and add initial data if not already present
func createFacilityCollection(pctx context.Context, db *mongo.Database, facilityName string) error {
	facilityCollection := db.Collection("facility")

	// Check if facility already exists in the database
	var existingFacility facility.Facilitiy
	err := facilityCollection.FindOne(pctx, bson.M{"name": facilityName}).Decode(&existingFacility)
	if err == nil {
		// Facility already exists, so skip insertion
		fmt.Printf("Facility %s already exists, skipping initialization.\n", facilityName)
		return nil
	} else if err != mongo.ErrNoDocuments {
		// Return any unexpected error
		return err
	}

	// Facility does not exist, so insert initial data
	var priceInsider, priceOutsider float64
	switch facilityName {
	case "fitness_facility":
		priceInsider, priceOutsider = 30.0, 40.0
	case "swimming_facility":
		priceInsider, priceOutsider = 40.0, 80.0
	case "football_facility":
		priceInsider, priceOutsider = 300.0, 400.0
	case "badminton_facility":
		priceInsider, priceOutsider = 80.0, 120.0
	}

	initialFacility := facility.Facilitiy{
		Id:           primitive.NewObjectID(),
		Name:         facilityName,
		PriceInsider: priceInsider,
		PriceOutsider: priceOutsider,
		Description:  fmt.Sprintf("Description of %s", facilityName),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	_, err = facilityCollection.InsertOne(pctx, initialFacility)
	return err
}

// Convert slots to []interface{}
func convertToInterface(slots []facility.Slot) []interface{} {
	interfaceSlots := make([]interface{}, len(slots))
	for i, slot := range slots {
		interfaceSlots[i] = slot
	}
	return interfaceSlots
}

// Create normal slots for fitness, swimming, and football facilities
func createNormalSlots(pctx context.Context, db *mongo.Database, facilityName string) error {
	slotCollection := db.Collection("slot")
	slots := []facility.Slot{}

	switch facilityName {
	case "fitness_facility":
		// Slots from 9:00 to 19:00 with 15 mins cleanup time
		startTime := time.Date(0, 1, 1, 9, 0, 0, 0, time.UTC)
		endTime := time.Date(0, 1, 1, 19, 0, 0, 0, time.UTC)
		for start := startTime; start.Before(endTime); start = start.Add(2*time.Hour + 15*time.Minute) {
			slot := facility.Slot{
				Id:             primitive.NewObjectID(),
				StartTime:      start.Format("15:04"),
				EndTime:        start.Add(2 * time.Hour).Format("15:04"),
				Status:         0,
				MaxBookings:    10,
				CurrentBookings: 0,
				FacilityType:   "fitness",
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			}
			slots = append(slots, slot)
		}
	case "swimming_facility", "football_facility":
		// Slots from 9:00 to 19:00 without cleanup time
		startTime := time.Date(0, 1, 1, 9, 0, 0, 0, time.UTC)
		endTime := time.Date(0, 1, 1, 19, 0, 0, 0, time.UTC)
		maxBookings := 1
		if facilityName == "swimming_facility" {

			maxBookings = 10
		}
		for start := startTime; start.Before(endTime); start = start.Add(2 * time.Hour) {
			slot := facility.Slot{
				Id:             primitive.NewObjectID(),
				StartTime:      start.Format("15:04"),
				EndTime:        start.Add(2 * time.Hour).Format("15:04"),
				Status:         0,
				MaxBookings:    maxBookings,
				CurrentBookings: 0,
				FacilityType:   facilityName,
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			}
			slots = append(slots, slot)
		}
	}

	// Insert slots as []interface{}
	_, err := slotCollection.InsertMany(pctx, convertToInterface(slots))
	return err
}

// Convert badminton slots to []interface{}
func convertBadmintonToInterface(slots []facility.BadmintonSlot) []interface{} {
	interfaceSlots := make([]interface{}, len(slots))
	for i, slot := range slots {
		interfaceSlots[i] = slot
	}
	return interfaceSlots
}

// Create courts for badminton facility
func createBadmintonCourts(pctx context.Context, db *mongo.Database) error {
	courtCollection := db.Collection("court")
	courts := []interface{}{
		facility.BadmintonCourt{Id: primitive.NewObjectID(), CourtNumber: 1, Status: 0},
		facility.BadmintonCourt{Id: primitive.NewObjectID(), CourtNumber: 2, Status: 0},
		facility.BadmintonCourt{Id: primitive.NewObjectID(), CourtNumber: 3, Status: 0},
		facility.BadmintonCourt{Id: primitive.NewObjectID(), CourtNumber: 4, Status: 0},
	}

	_, err := courtCollection.InsertMany(pctx, courts)
	return err
}

// Create badminton slots from 10:00 to 21:00
func createBadmintonSlots(pctx context.Context, db *mongo.Database) error {
	slotCollection := db.Collection("badminton_slot")
	slots := []facility.BadmintonSlot{}
	startTime := time.Date(0, 1, 1, 10, 0, 0, 0, time.UTC)
	endTime := time.Date(0, 1, 1, 21, 0, 0, 0, time.UTC)

	for start := startTime; start.Before(endTime); start = start.Add(2 * time.Hour) {
		for courtNumber := 1; courtNumber <= 4; courtNumber++ {
			slot := facility.BadmintonSlot{
				Id:             primitive.NewObjectID(),
				StartTime:      start.Format("15:04"),
				EndTime:        start.Add(2 * time.Hour).Format("15:04"),
				CourtId:        primitive.NewObjectID(), // Replace with actual court ID in a full setup
				Status:         0,
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			}
			slots = append(slots, slot)
		}
	}

	// Insert slots as []interface{}
	_, err := slotCollection.InsertMany(pctx, convertBadmintonToInterface(slots))
	return err
}
