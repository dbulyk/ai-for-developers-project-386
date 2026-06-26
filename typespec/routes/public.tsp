import "@typespec/http";

using Http;

namespace CalendarBooking;

/** Public endpoints accessible by guests without authentication. */
@route("/public")
namespace Public {
  /** Browse available event types. */
  @route("/event-types")
  interface PublicEventTypes {
    /** List all event types with their name, description, and duration. */
    @get
    list(): EventType[];

    /** Get available booking slots for an event type.
        Returns free slots for the next 14 days (Mon–Fri, 09:00–18:00 owner local time)
        grouped by calendar day. Slot grid matches the event type duration. */
    @get
    @route("/{id}/slots")
    getSlots(@path id: string): AvailableDay[] | NotFoundResponse;
  }

  /** Create guest bookings. */
  @route("/bookings")
  interface PublicBookings {
    /** Book a slot. Returns 409 if the requested slot is already taken. */
    @post
    create(@body body: CreateBookingRequest): {
      @statusCode statusCode: 201;
      @body body: Booking;
    } | NotFoundResponse | ConflictResponse;
  }
}
