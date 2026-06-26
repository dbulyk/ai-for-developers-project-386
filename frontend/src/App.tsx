import { Routes, Route, Navigate } from 'react-router-dom';
import { EventTypesPage } from './pages/EventTypesPage';
import { BookingPage } from './pages/BookingPage';
import { AdminLayout } from './pages/admin/AdminLayout';
import { AdminEventTypesPage } from './pages/admin/AdminEventTypesPage';
import { AdminBookingsPage } from './pages/admin/AdminBookingsPage';

export function App() {
  return (
    <Routes>
      <Route path="/" element={<EventTypesPage />} />
      <Route path="/event-types/:id" element={<BookingPage />} />
      <Route path="/admin" element={<AdminLayout />}>
        <Route index element={<Navigate to="/admin/event-types" replace />} />
        <Route path="event-types" element={<AdminEventTypesPage />} />
        <Route path="bookings" element={<AdminBookingsPage />} />
      </Route>
    </Routes>
  );
}
