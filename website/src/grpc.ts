import { Timestamp } from 'google-protobuf/google/protobuf/timestamp_pb';
import { Events } from './sdk/events_pb';

const host = window.location.origin + '/api';

export const grpc = {
  events: new Events(host),
}

// https://github.com/SafetyCulture/grpc-web-devtools
const devtools = (window as any).__GRPCWEB_DEVTOOLS__;
if (devtools) {
  devtools(Object.values(grpc));
}

export function toDate(timestamp: Timestamp.AsObject): Date {
  const t = new Timestamp();
  t.setSeconds(timestamp.seconds);
  t.setNanos(timestamp.nanos);
  return t.toDate();
}

export function toTimestamp(date: Date): Timestamp.AsObject {
  const t = new Timestamp();
  t.fromDate(date);
  return t.toObject();
}
