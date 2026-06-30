export function exactText(name: string): RegExp {
  return new RegExp('^' + name.replace(/[.*+?^${}()|[\]\\]/g, '\\$&') + '$');
}

export const Selectors = {
  // Event types page
  eventTypesTitle: 'event-types-title',
  eventTypeCard: 'event-type-card',
  eventTypeName: 'event-type-name',
  eventTypeDuration: 'event-type-duration',
  eventTypeBookButton: 'event-type-book-button',

  // Booking page
  bookingPageTitle: 'booking-page-title',
  bookingPageDuration: 'booking-page-duration',
  dayButton: 'day-button',
  slotCard: 'slot-card',
  slotStatusAttr: 'data-status',
  bookingModal: 'booking-modal',
  guestNameInput: 'guest-name-input',
  confirmBookingButton: 'confirm-booking-button',
  cancelBookingModalButton: 'cancel-booking-modal-button',

  // Admin event types page
  newEventTypeButton: 'new-event-type-button',
  eventTypesTable: 'event-types-table',
  eventTypeRow: 'event-type-row',
  eventTypeRowName: 'event-type-row-name',
  editEventTypeButton: 'edit-event-type-button',
  deleteEventTypeButton: 'delete-event-type-button',
  createEventTypeModal: 'create-event-type-modal',
  editEventTypeModal: 'edit-event-type-modal',
  deleteEventTypeModal: 'delete-event-type-modal',
  eventTypeNameInput: 'event-type-name-input',
  eventTypeDescriptionInput: 'event-type-description-input',
  eventTypeDurationInput: 'event-type-duration-input',
  createEventTypeSubmit: 'create-event-type-submit',
  saveEventTypeSubmit: 'save-event-type-submit',
  confirmDeleteEventType: 'confirm-delete-event-type',

  // Admin bookings page
  bookingsTitle: 'bookings-title',
  bookingsTable: 'bookings-table',
  bookingRow: 'booking-row',
  bookingRowGuest: 'booking-row-guest',
  cancelBookingButton: 'cancel-booking-button',
  cancelBookingModal: 'cancel-booking-modal',
  confirmCancelBooking: 'confirm-cancel-booking',
} as const;
