package configs

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/supabase-community/supabase-go"
)

var SupabaseClient *supabase.Client

// InitializeSupabase initializes the Supabase client
func InitializeSupabase() error {
	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_ANON_KEY")

	if supabaseURL == "" || supabaseKey == "" {
		return fmt.Errorf("SUPABASE_URL and SUPABASE_ANON_KEY environment variables are required")
	}

	client, err := supabase.NewClient(supabaseURL, supabaseKey, nil)
	if err != nil {
		return fmt.Errorf("failed to create Supabase client: %w", err)
	}

	// Test the connection
	ctx := context.Background()
	_, err = client.DB.From("users").Select("count", false, "", "", "").Execute(ctx)
	if err != nil {
		return fmt.Errorf("failed to connect to Supabase: %w", err)
	}

	SupabaseClient = client
	log.Println("✅ Supabase connection established successfully")
	return nil
}

// GetSupabaseClient returns the initialized Supabase client
func GetSupabaseClient() *supabase.Client {
	return SupabaseClient
}

// CloseSupabase closes the Supabase connection
func CloseSupabase() error {
	if SupabaseClient != nil {
		// Supabase client doesn't have a close method, but we can clean up resources
		log.Println("Supabase connection closed")
	}
	return nil
} 