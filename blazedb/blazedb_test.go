package blazedb

import (
	"testing"

	bolt "go.etcd.io/bbolt"
)

func TestNew(t *testing.T) {
	// Create a temporary database for testing
	db, err := bolt.Open("test.db", 0600, nil)
	if err != nil {
		t.Fatalf("Failed to open temporary database: %v", err)
	}
	defer db.Close()
	defer func() {
		// Cleanup after testing
		err := db.Close()
		if err != nil {
			t.Fatalf("Failed to close temporary database: %v", err)
		}
	}()

	// Define your test cases
	testCases := []struct {
		name        string
		options     []OptFunc
		expectedDB  *BlazeDB
		expectedErr error
	}{
		{
			name:    "DefaultOptions",
			options: nil,
			expectedDB: &BlazeDB{
				currentDatabase: "defaultDBName.blaze",
				Options:         &Options{DBName: "defaultDBName"},
				db:              db,
			},
			expectedErr: nil,
		},
		// Add more test cases for other scenarios and edge cases as needed
	}

	// Run the test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, err := New(tc.options...)
			if err != tc.expectedErr {
				t.Errorf("Expected error: %v, got: %v", tc.expectedErr, err)
				return
			}
			if err == nil && db != nil {
				// Check if the returned database has the expected fields
				if db.currentDatabase != tc.expectedDB.currentDatabase {
					t.Errorf("Expected currentDatabase: %s, got: %s", tc.expectedDB.currentDatabase, db.currentDatabase)
				}
				if db.Options.DBName != tc.expectedDB.Options.DBName {
					t.Errorf("Expected Options.DBName: %s, got: %s", tc.expectedDB.Options.DBName, db.Options.DBName)
				}
				// Add more checks if needed
			}
		})
	}
}
