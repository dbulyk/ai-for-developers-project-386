namespace CalendarBooking;

/** A confirmed booking for a specific event type slot. */
model Booking {
  id: string;
  eventTypeId: string;

  @minLength(1)
  guestName: string;

  /** Slot start time in UTC. */
  startTime: utcDateTime;

  /** Slot end time in UTC, computed as startTime + durationMinutes. */
  endTime: utcDateTime;

  createdAt: utcDateTime;
}

/** Request body for creating a new booking. */
model CreateBookingRequest {
  eventTypeId: string;

  @minLength(1)
  guestName: string;

  /** Desired slot start time in UTC. Must be a free slot within the 14-day window. */
  startTime: utcDateTime;
}
