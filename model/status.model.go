package model

/**
 * DB Info
 */

/**
 * SEARCH regex fields
 */

/**
 * MODEL
 */
type Status string

/**
 * ENUM
 */
var (
	// StatusActive ...
	StatusActive Status = "ACTIVE"
	// StatusDraft ...
	StatusDraft Status = "DRAFT"
	// StatusArchive ...
	StatusArchive Status = "ARCHIVE"
)
