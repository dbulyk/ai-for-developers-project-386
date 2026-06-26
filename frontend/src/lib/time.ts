import dayjs from 'dayjs';
import utc from 'dayjs/plugin/utc';
import timezone from 'dayjs/plugin/timezone';

dayjs.extend(utc);
dayjs.extend(timezone);

export const ownerTz = (): string =>
  (import.meta.env.VITE_OWNER_TIMEZONE as string | undefined) ?? 'Europe/Moscow';

export function formatSlotTime(utcIso: string): string {
  return dayjs.utc(utcIso).tz(ownerTz()).format('HH:mm');
}

export function formatDateTime(utcIso: string): string {
  return dayjs.utc(utcIso).tz(ownerTz()).format('DD MMM YYYY, HH:mm');
}
