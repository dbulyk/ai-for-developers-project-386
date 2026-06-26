import "@typespec/http";
import "@typespec/openapi3";

import "./models/event-type.tsp";
import "./models/booking.tsp";
import "./models/slot.tsp";
import "./routes/admin.tsp";
import "./routes/public.tsp";

using Http;

/** REST API for calendar-based event booking. */
@service(#{ title: "Calendar Booking API" })
@server("http://localhost:8080", "Development server")
namespace CalendarBooking;
