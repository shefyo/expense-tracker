package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strconv"
	"time"
)

type Expense struct {
	ID          int       `json:"id"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
	Category    string    `json:"category"`
}

type Budget struct {
	Month  int     `json:"month"`
	Amount float64 `json:"amount"`
}

type Store struct {
	Expenses []Expense `json:"expenses"`
	Budgets  []Budget  `json:"budgets"`
	NextID   int       `json:"next_id"`
}

const dataFile = "expenses.json"

func loadStore() Store {
	store := Store{NextID: 1}
	data, err := os.ReadFile(dataFile)
	if err != nil {
		return store
	}
	json.Unmarshal(data, &store)
	return store
}

func saveStore(store Store) {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		fmt.Println("Error saving data:", err)
		os.Exit(1)
	}
	if err := os.WriteFile(dataFile, data, 0644); err != nil {
		fmt.Println("Error writing file:", err)
		os.Exit(1)
	}
}

func parseFlags(args []string) map[string]string {
	flags := make(map[string]string)
	for i := 0; i < len(args); i++ {
		if len(args[i]) > 2 && args[i][:2] == "--" {
			key := args[i][2:]
			if i+1 < len(args) && (len(args[i+1]) < 2 || args[i+1][:2] != "--") {
				flags[key] = args[i+1]
				i++
			} else {
				flags[key] = "true"
			}
		}
	}
	return flags
}

func monthName(m int) string {
	months := []string{"January", "February", "March", "April", "May", "June",
		"July", "August", "September", "October", "November", "December"}
	if m < 1 || m > 12 {
		return "Unknown"
	}
	return months[m-1]
}

func cmdAdd(flags map[string]string) {
	desc, ok := flags["description"]
	if !ok || desc == "" {
		fmt.Println("Error: --description is required")
		os.Exit(1)
	}
	amtStr, ok := flags["amount"]
	if !ok {
		fmt.Println("Error: --amount is required")
		os.Exit(1)
	}
	amt, err := strconv.ParseFloat(amtStr, 64)
	if err != nil || amt <= 0 {
		fmt.Println("Error: --amount must be a positive number")
		os.Exit(1)
	}
	amt = math.Round(amt*100) / 100

	category := flags["category"]

	store := loadStore()
	exp := Expense{
		ID:          store.NextID,
		Date:        time.Now(),
		Description: desc,
		Amount:      amt,
		Category:    category,
	}
	store.Expenses = append(store.Expenses, exp)
	store.NextID++
	saveStore(store)
	fmt.Printf("Expense added successfully (ID: %d)\n", exp.ID)
}

func cmdUpdate(flags map[string]string) {
	idStr, ok := flags["id"]
	if !ok {
		fmt.Println("Error: --id is required")
		os.Exit(1)
	}
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		fmt.Println("Error: --id must be a positive integer")
		os.Exit(1)
	}

	store := loadStore()
	found := false
	for i, exp := range store.Expenses {
		if exp.ID == id {
			found = true
			if desc, ok := flags["description"]; ok && desc != "" {
				store.Expenses[i].Description = desc
			}
			if amtStr, ok := flags["amount"]; ok {
				amt, err := strconv.ParseFloat(amtStr, 64)
				if err != nil || amt <= 0 {
					fmt.Println("Error: --amount must be a positive number")
					os.Exit(1)
				}
				store.Expenses[i].Amount = math.Round(amt*100) / 100
			}
			if cat, ok := flags["category"]; ok {
				store.Expenses[i].Category = cat
			}
			break
		}
	}
	if !found {
		fmt.Printf("Error: expense with ID %d not found\n", id)
		os.Exit(1)
	}
	saveStore(store)
	fmt.Printf("Expense %d updated successfully\n", id)
}

func cmdDelete(flags map[string]string) {
	idStr, ok := flags["id"]
	if !ok {
		fmt.Println("Error: --id is required")
		os.Exit(1)
	}
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		fmt.Println("Error: --id must be a positive integer")
		os.Exit(1)
	}

	store := loadStore()
	newExpenses := make([]Expense, 0, len(store.Expenses))
	found := false
	for _, exp := range store.Expenses {
		if exp.ID == id {
			found = true
			continue
		}
		newExpenses = append(newExpenses, exp)
	}
	if !found {
		fmt.Printf("Error: expense with ID %d not found\n", id)
		os.Exit(1)
	}
	store.Expenses = newExpenses
	saveStore(store)
	fmt.Println("Expense deleted successfully")
}

func cmdList(flags map[string]string) {
	store := loadStore()
	expenses := store.Expenses

	if cat, ok := flags["category"]; ok && cat != "" {
		filtered := expenses[:0]
		for _, e := range expenses {
			if e.Category == cat {
				filtered = append(filtered, e)
			}
		}
		expenses = filtered
	}

	if len(expenses) == 0 {
		fmt.Println("No expenses found.")
		return
	}

	fmt.Printf("%-4s  %-12s  %-20s  %-10s  %s\n", "ID", "Date", "Description", "Amount", "Category")
	for _, e := range expenses {
		cat := e.Category
		if cat == "" {
			cat = "-"
		}
		fmt.Printf("%-4d  %-12s  %-20s  $%-9.2f  %s\n",
			e.ID,
			e.Date.Format("2006-01-02"),
			e.Description,
			e.Amount,
			cat,
		)
	}
}

func cmdSummary(flags map[string]string) {
	store := loadStore()
	expenses := store.Expenses

	if monthStr, ok := flags["month"]; ok {
		month, err := strconv.Atoi(monthStr)
		if err != nil || month < 1 || month > 12 {
			fmt.Println("Error: --month must be between 1 and 12")
			os.Exit(1)
		}
		currentYear := time.Now().Year()
		var total float64
		for _, e := range expenses {
			if int(e.Date.Month()) == month && e.Date.Year() == currentYear {
				total += e.Amount
			}
		}

		budget := -1.0
		for _, b := range store.Budgets {
			if b.Month == month {
				budget = b.Amount
				break
			}
		}

		fmt.Printf("Total expenses for %s: $%.2f\n", monthName(month), total)
		if budget >= 0 {
			fmt.Printf("Budget for %s: $%.2f\n", monthName(month), budget)
			if total > budget {
				fmt.Printf("⚠  Warning: you have exceeded your budget by $%.2f\n", total-budget)
			} else {
				fmt.Printf("Remaining budget: $%.2f\n", budget-total)
			}
		}
		return
	}

	var total float64
	for _, e := range expenses {
		total += e.Amount
	}
	fmt.Printf("Total expenses: $%.2f\n", total)
}

func cmdBudget(flags map[string]string) {
	monthStr, ok := flags["month"]
	if !ok {
		fmt.Println("Error: --month is required")
		os.Exit(1)
	}
	month, err := strconv.Atoi(monthStr)
	if err != nil || month < 1 || month > 12 {
		fmt.Println("Error: --month must be between 1 and 12")
		os.Exit(1)
	}
	amtStr, ok := flags["amount"]
	if !ok {
		fmt.Println("Error: --amount is required")
		os.Exit(1)
	}
	amt, err := strconv.ParseFloat(amtStr, 64)
	if err != nil || amt <= 0 {
		fmt.Println("Error: --amount must be a positive number")
		os.Exit(1)
	}

	store := loadStore()
	found := false
	for i, b := range store.Budgets {
		if b.Month == month {
			store.Budgets[i].Amount = amt
			found = true
			break
		}
	}
	if !found {
		store.Budgets = append(store.Budgets, Budget{Month: month, Amount: amt})
	}
	saveStore(store)
	fmt.Printf("Budget for %s set to $%.2f\n", monthName(month), amt)
}

func cmdExport(flags map[string]string) {
	outFile := flags["output"]
	if outFile == "" {
		outFile = "expenses.csv"
	}

	store := loadStore()
	f, err := os.Create(outFile)
	if err != nil {
		fmt.Println("Error creating file:", err)
		os.Exit(1)
	}
	defer f.Close()

	f.WriteString("ID,Date,Description,Amount,Category\n")
	for _, e := range store.Expenses {
		cat := e.Category
		if cat == "" {
			cat = "-"
		}
		line := fmt.Sprintf("%d,%s,\"%s\",%.2f,%s\n",
			e.ID,
			e.Date.Format("2006-01-02"),
			e.Description,
			e.Amount,
			cat,
		)
		f.WriteString(line)
	}
	fmt.Printf("Expenses exported to %s\n", outFile)
}

func printUsage() {
	fmt.Println(`Usage: expense-tracker <command> [options]

Commands:
  add        Add a new expense
             --description  Description of the expense (required)
             --amount       Amount of the expense (required)
             --category     Category of the expense (optional)

  update     Update an existing expense
             --id           ID of the expense to update (required)
             --description  New description (optional)
             --amount       New amount (optional)
             --category     New category (optional)

  delete     Delete an expense
             --id           ID of the expense to delete (required)

  list       List all expenses
             --category     Filter by category (optional)

  summary    Show expense summary
             --month        Show summary for a specific month (1-12, optional)

  budget     Set a monthly budget
             --month        Month number 1-12 (required)
             --amount       Budget amount (required)

  export     Export expenses to CSV
             --output       Output file path (default: expenses.csv)`)
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]
	flags := parseFlags(os.Args[2:])

	switch command {
	case "add":
		cmdAdd(flags)
	case "update":
		cmdUpdate(flags)
	case "delete":
		cmdDelete(flags)
	case "list":
		cmdList(flags)
	case "summary":
		cmdSummary(flags)
	case "budget":
		cmdBudget(flags)
	case "export":
		cmdExport(flags)
	case "help", "--help", "-h":
		printUsage()
	default:
		fmt.Printf("Error: unknown command %q\n\n", command)
		printUsage()
		os.Exit(1)
	}
}
