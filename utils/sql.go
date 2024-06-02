package utils

import "database/sql"

func SqlNullStringIfValidFromString[T ~string](s T) sql.NullString {
	if s == "" {
		return sql.NullString{}
	}
	return sql.NullString{Valid: true, String: string(s)}
}

func SqlNullStringIfValidFromStringPtr[T ~*string](s T) sql.NullString {
	if s == nil {
		return sql.NullString{}
	}
	return sql.NullString{Valid: true, String: *s}
}
