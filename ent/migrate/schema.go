// Code generated by entc, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// CheckinsColumns holds the columns for the "checkins" table.
	CheckinsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "checkin_time", Type: field.TypeInt64},
		{Name: "event_id", Type: field.TypeString},
		{Name: "user_checkins", Type: field.TypeInt, Nullable: true},
	}
	// CheckinsTable holds the schema information for the "checkins" table.
	CheckinsTable = &schema.Table{
		Name:       "checkins",
		Columns:    CheckinsColumns,
		PrimaryKey: []*schema.Column{CheckinsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "checkins_users_checkins",
				Columns:    []*schema.Column{CheckinsColumns[3]},
				RefColumns: []*schema.Column{UsersColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// UsersColumns holds the columns for the "users" table.
	UsersColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "phone_number", Type: field.TypeString},
	}
	// UsersTable holds the schema information for the "users" table.
	UsersTable = &schema.Table{
		Name:       "users",
		Columns:    UsersColumns,
		PrimaryKey: []*schema.Column{UsersColumns[0]},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		CheckinsTable,
		UsersTable,
	}
)

func init() {
	CheckinsTable.ForeignKeys[0].RefTable = UsersTable
}
