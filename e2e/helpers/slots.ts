export const OWNER_TIMEZONE = 'Europe/Moscow';

export interface Slot {
  startTime: string;
  status: 'free' | 'taken';
}

export interface AvailableDay {
  date: string;
  slots: Slot[];
}

export function firstFreeSlot(days: AvailableDay[]): Slot | undefined {
  for (const day of days) {
    const slot = day.slots.find((s) => s.status === 'free');
    if (slot) return slot;
  }
  return undefined;
}

export function firstTakenSlot(days: AvailableDay[]): Slot | undefined {
  for (const day of days) {
    const slot = day.slots.find((s) => s.status === 'taken');
    if (slot) return slot;
  }
  return undefined;
}

export function formatSlotTime(utcIso: string): string {
  return new Date(utcIso).toLocaleTimeString('ru-RU', {
    timeZone: OWNER_TIMEZONE,
    hour: '2-digit',
    minute: '2-digit',
    hour12: false,
  });
}

export function formatDateTime(utcIso: string): string {
  return new Date(utcIso).toLocaleString('en-GB', {
    timeZone: OWNER_TIMEZONE,
    day: '2-digit',
    month: 'short',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
    hour12: false,
  });
}

export function formatDayLabel(date: string): string {
  return date;
}
