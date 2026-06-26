import "@typespec/http";

using Http;

namespace CalendarBooking;

/** Admin endpoints for the calendar owner. No authentication required — access is
    restricted at the infrastructure level to the pre-configured owner profile. */
@route("/admin")
namespace Admin {
  /** Manage event types. */
  @route("/event-types")
  interface AdminEventTypes {
    /** List all event types. */
    @get
    list(): EventType[];

    /** Create a new event type. */
    @post
    create(@body body: CreateEventTypeRequest): {
      @statusCode statusCode: 201;
      @body body: EventType;
    };

    /** Replace an existing event type (full update). */
    @put
    @route("/{id}")
    update(@path id: string, @body body: UpdateEventTypeRequest): EventType | NotFoundResponse;

    /** Delete an event type. */
    @delete
    @route("/{id}")
    delete(@path id: string): NoContentResponse | NotFoundResponse;
  }

  /** Manage bookings made by guests. */
  @route("/bookings")
  interface AdminBookings {
    /** List all upcoming bookings across all event types, ordered by start time. */
    @get
    list(): Booking[];

    /** Cancel a booking. */
    @delete
    @route("/{id}")
    cancel(@path id: string): NoContentResponse | NotFoundResponse;
  }
}
