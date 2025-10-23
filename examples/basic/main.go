package main

import (
	"fmt"
	"net/url"

	"github.com/lucasvillarinho/restql"
)

func main() {
	fmt.Println("=== RestQL Examples ===")
	fmt.Println()

	// Create schema for users table
	usersSchema := restql.NewSchema("users").
		AllowFields("id", "name", "email", "status", "age", "created", "created_at", "verified")

	employeesSchema := restql.NewSchema("employees").
		AllowFields("id", "name", "age", "salary")

	ordersSchema := restql.NewSchema("orders").
		AllowFields("id", "status", "total")

	postsSchema := restql.NewSchema("posts").
		AllowFields("id", "title", "deleted_at", "approved_at")

	// Query 1: Simple equality filter
	fmt.Println("Query 1: Simple equality filter")
	params1, _ := url.ParseQuery("filter=status='active'")
	qb1, _ := restql.Parse(params1, usersSchema)
	sql1, args1 := qb1.ToSQL()
	fmt.Printf("SQL: %s\n", sql1)
	fmt.Printf("Args: %v\n\n", args1)
	// Output:
	// SQL: SELECT * FROM users WHERE status = ?
	// Args: [active]

	// Query 2: Multiple conditions with AND
	fmt.Println("Query 2: Multiple conditions with AND")
	params2, _ := url.ParseQuery("filter=status='active' && age>=18")
	qb2, _ := restql.Parse(params2, usersSchema)
	sql2, args2 := qb2.ToSQL()
	fmt.Printf("SQL: %s\n", sql2)
	fmt.Printf("Args: %v\n\n", args2)
	// Output:
	// SQL: SELECT * FROM users WHERE (status = ? AND age >= ?)
	// Args: [active 18]

	// Query 3: OR conditions
	fmt.Println("Query 3: OR conditions")
	params3, _ := url.ParseQuery("filter=status='active' || status='pending'")
	qb3, _ := restql.Parse(params3, usersSchema)
	sql3, args3 := qb3.ToSQL()
	fmt.Printf("SQL: %s\n", sql3)
	fmt.Printf("Args: %v\n\n", args3)
	// Output:
	// SQL: SELECT * FROM users WHERE (status = ? OR status = ?)
	// Args: [active pending]

	// Query 4: Comparison operators
	fmt.Println("Query 4: Comparison operators")
	params4, _ := url.ParseQuery("filter=age>18 && salary<=50000")
	qb4, _ := restql.Parse(params4, employeesSchema)
	sql4, args4 := qb4.ToSQL()
	fmt.Printf("SQL: %s\n", sql4)
	fmt.Printf("Args: %v\n\n", args4)
	// Output:
	// SQL: SELECT * FROM employees WHERE (age > ? AND salary <= ?)
	// Args: [18 50000]

	// Query 5: LIKE operator (~)
	fmt.Println("Query 5: LIKE operator (~)")
	params5, _ := url.ParseQuery("filter=name~'John' && email~'@gmail.com'")
	qb5, _ := restql.Parse(params5, usersSchema)
	sql5, args5 := qb5.ToSQL()
	fmt.Printf("SQL: %s\n", sql5)
	fmt.Printf("Args: %v\n\n", args5)
	// Output:
	// SQL: SELECT * FROM users WHERE (name LIKE ? AND email LIKE ?)
	// Args: [John @gmail.com]

	// Query 6: IN operator
	fmt.Println("Query 6: IN operator")
	params6, _ := url.ParseQuery("filter=status IN ['active','pending','approved']")
	qb6, _ := restql.Parse(params6, ordersSchema)
	sql6, args6 := qb6.ToSQL()
	fmt.Printf("SQL: %s\n", sql6)
	fmt.Printf("Args: %v\n\n", args6)
	// Output:
	// SQL: SELECT * FROM orders WHERE status IN (?, ?, ?)
	// Args: [active pending approved]

	// Query 7: IS NULL / IS NOT NULL
	fmt.Println("Query 7: IS NULL / IS NOT NULL")
	params7, _ := url.ParseQuery("filter=deleted_at IS NULL && approved_at IS NOT NULL")
	qb7, _ := restql.Parse(params7, postsSchema)
	sql7, args7 := qb7.ToSQL()
	fmt.Printf("SQL: %s\n", sql7)
	fmt.Printf("Args: %v\n\n", args7)
	// Output:
	// SQL: SELECT * FROM posts WHERE (deleted_at IS NULL AND approved_at IS NOT NULL)
	// Args: []

	// Query 8: Complex nested conditions
	fmt.Println("Query 8: Complex nested conditions")
	params8, _ := url.ParseQuery("filter=(status='active' && age>=18) || (status='premium' && age>=16)")
	qb8, _ := restql.Parse(params8, usersSchema)
	sql8, args8 := qb8.ToSQL()
	fmt.Printf("SQL: %s\n", sql8)
	fmt.Printf("Args: %v\n\n", args8)
	// Output:
	// SQL: SELECT * FROM users WHERE ((status = ? AND age >= ?) OR (status = ? AND age >= ?))
	// Args: [active 18 premium 16]

	// Query 9: Sort (ORDER BY)
	fmt.Println("Query 9: Sort (ORDER BY)")
	params9, _ := url.ParseQuery("filter=status='active'&sort=-created,name")
	qb9, _ := restql.Parse(params9, usersSchema)
	sql9, args9 := qb9.ToSQL()
	fmt.Printf("SQL: %s\n", sql9)
	fmt.Printf("Args: %v\n\n", args9)
	// Output:
	// SQL: SELECT * FROM users WHERE status = ? ORDER BY created DESC, name ASC
	// Args: [active]

	// Query 10: Limit and Offset (pagination)
	fmt.Println("Query 10: Limit and Offset (pagination)")
	params10, _ := url.ParseQuery("filter=status='active'&limit=10&offset=20")
	qb10, _ := restql.Parse(params10, usersSchema)
	sql10, args10 := qb10.ToSQL()
	fmt.Printf("SQL: %s\n", sql10)
	fmt.Printf("Args: %v\n\n", args10)
	// Output:
	// SQL: SELECT * FROM users WHERE status = ? LIMIT 10 OFFSET 20
	// Args: [active]

	// Query 11: Select specific fields
	fmt.Println("Query 11: Select specific fields")
	params11, _ := url.ParseQuery("fields=id,name,email&filter=status='active'")
	qb11, _ := restql.Parse(params11, usersSchema)
	sql11, args11 := qb11.ToSQL()
	fmt.Printf("SQL: %s\n", sql11)
	fmt.Printf("Args: %v\n\n", args11)
	// Output:
	// SQL: SELECT id, name, email FROM users WHERE status = ?
	// Args: [active]

	// Query 12: Full example with all features
	fmt.Println("Query 12: Full example with all features")
	params12, _ := url.ParseQuery("fields=id,name,email,status,created_at&filter=(status='active' && age>=18) || (status='premium' && verified=true)&sort=-created_at,name&limit=20&offset=0")
	qb12, _ := restql.Parse(params12, usersSchema)
	sql12, args12 := qb12.ToSQL()
	fmt.Printf("SQL: %s\n", sql12)
	fmt.Printf("Args: %v\n\n", args12)
	// Output:
	// SQL: SELECT id, name, email, status, created_at FROM users WHERE ((status = ? AND age >= ?) OR (status = ? AND verified = ?)) ORDER BY created_at DESC, name ASC LIMIT 20 OFFSET 0
	// Args: [active 18 premium true]
}
