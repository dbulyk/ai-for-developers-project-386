namespace CalendarBooking;

/** Status of a generated slot. */
enum SlotStatus {
  Free: "free",
  Taken: "taken",
}

/** A single booking slot with its availability status. */
model Slot {
  /** Slot start time in UTC. */
  startTime: utcDateTime;

  /** Whether the slot is free or already taken. */
  status: SlotStatus;
}

/** Available booking slots for a single calendar day. */
model AvailableDay {
  /** Calendar date in ISO 8601 format (e.g. "2026-06-27"). */
  date: plainDate;

  /** Slots for this day in chronological order, including taken ones. */
  slots: Slot[];
}
