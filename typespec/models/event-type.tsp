namespace CalendarBooking;

/** A type of bookable event defined by the calendar owner. */
model EventType {
  id: string;

  @minLength(1)
  name: string;

  description: string;

  /** Duration of the event in minutes. Determines slot grid size. */
  @minValue(1)
  durationMinutes: int32;
}

/** Request body for creating a new event type. */
model CreateEventTypeRequest {
  @minLength(1)
  name: string;

  description: string;

  @minValue(1)
  durationMinutes: int32;
}

/** Request body for replacing an existing event type (full update). */
model UpdateEventTypeRequest {
  @minLength(1)
  name: string;

  description: string;

  @minValue(1)
  durationMinutes: int32;
}
