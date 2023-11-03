package main

import (
	"database/sql"
	"fmt"
	"log"
	"regexp"

	_ "github.com/mattn/go-sqlite3"
)

var (
	originalPhoneNumbers = []string{
		"1234567890",
		"123 456 7891",
		"(123) 456 7892",
		"(123) 456-7893",
		"123-456-7894",
		"123-456-7890",
		"1234567892",
		"(123)456-7892",
	}
)

func main() {
	db, err := sql.Open("sqlite3", "phone.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	defer db.Close()

	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS phone_numbers (phone_number TEXT)`); err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	if _, err := db.Exec(`DELETE FROM phone_numbers`); err != nil {
		log.Fatalf("Failed to delete rows: %v", err)
	}

	// for insert data from user, always user prepare for proper escaping to prevent SQLi vuln
	stmt, err := db.Prepare("INSERT INTO phone_numbers (phone_number) VALUES (?)")
	if err != nil {
		log.Fatalf("Failed to prepare statement: %v", err)
	}

	for _, phoneNumber := range originalPhoneNumbers {
		if _, err := stmt.Exec(phoneNumber); err != nil {
			log.Fatalf("Failed to insert row: %v", err)
		}
	}

	rows, err := db.Query("SELECT phone_number FROM phone_numbers")
	if err != nil {
		log.Fatalf("Failed to query rows: %v", err)
	}
	defer rows.Close()

	fmt.Println("Original data:")
	var originalPhoneNumbers []string
	var formatPhoneNumbers []string
	for rows.Next() {
		var phoneNumber string
		if err := rows.Scan(&phoneNumber); err != nil {
			log.Fatalf("Failed to scan row: %v", err)
		}

		formattedPhoneNumber := formatPhoneNumber(phoneNumber)
		originalPhoneNumbers = append(originalPhoneNumbers, phoneNumber)
		formatPhoneNumbers = append(formatPhoneNumbers, formattedPhoneNumber)

		fmt.Println(phoneNumber)
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("Failed to get rows error: %v", err)
	}

	stmtQuery, err := db.Prepare("SELECT phone_number FROM phone_numbers WHERE phone_number=? LIMIT 1")
	if err != nil {
		log.Fatalf("Failed to prepare statement: %v", err)
	}

	// Update data
	stmtUpdate, err := db.Prepare("UPDATE phone_numbers SET phone_number=? WHERE phone_number=?")
	if err != nil {
		log.Fatalf("Failed to prepare statement: %v", err)
	}

	stmtDelete, err := db.Prepare("DELETE FROM phone_numbers WHERE phone_number=?")
	if err != nil {
		log.Fatalf("Failed to prepare statement: %v", err)
	}

	for i, phoneNumber := range originalPhoneNumbers {
		formattedPhoneNumber := formatPhoneNumbers[i]
		res, err := stmtQuery.Query(formattedPhoneNumber)

		if err != nil {
			log.Fatalf("Failed to query row: %v", err)
		}

		duplicateFound := res.Next()
		res.Close()
		if duplicateFound {
			log.Printf("Duplicate found for %s, delete: %s", formattedPhoneNumber, phoneNumber)
			if _, err := stmtDelete.Exec(phoneNumber); err != nil {
				log.Fatalf("Failed to delete row: %v", err)
			}
			continue
		}

		log.Printf("Update: %s with %s", phoneNumber, formattedPhoneNumber)
		if _, err := stmtUpdate.Exec(formattedPhoneNumber, phoneNumber); err != nil {
			log.Fatalf("Failed to update row: %v", err)
		}
	}

	// List new data
	fmt.Println()
	fmt.Println("New Formatted data:")
	rows, err = db.Query("SELECT phone_number FROM phone_numbers ORDER BY phone_number")
	if err != nil {
		log.Fatalf("Failed to query rows: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var phoneNumber string
		if err := rows.Scan(&phoneNumber); err != nil {
			log.Fatalf("Failed to scan row: %v", err)
		}
		fmt.Println(phoneNumber)
	}
}

func formatPhoneNumber(phoneNumber string) string {
	// var formattedPhoneNumber string
	// // remove all non-digit characters
	// for _, char := range phoneNumber {
	// 	if char >= '0' && char <= '9' {
	// 		formattedPhoneNumber += string(char)
	// 	}
	// }
	// return phoneNumber

	// replace any thing non digit with empty string
	return regexp.MustCompile("[^\\d]").ReplaceAllString(phoneNumber, "")
}
