namespace CalendarBooking;

/** Available booking slots for a single calendar day. */
model AvailableDay {
  /** Calendar date in ISO 8601 format (e.g. "2026-06-27"). */
  date: plainDate;

  /** Free slot start times in UTC for this day. */
  slots: utcDateTime[];
}
