import {
  Radio,
  RadioGroup,
  Tag
} from "@blueprintjs/core";
import { subMinutes } from "date-fns";
import React, { useEffect, useState } from "react";
import { EventList } from "./EventList";
import { grpc, toTimestamp } from "./grpc";
import { LayoutContent } from "./Layout";
import { Event } from "./sdk/events_pb";

export function App() {
  const [events, setEvents] = useState(new Array<Event.AsObject>());
  const [since, setSince] = useState(30);

  useEffect(() => {
    const stream = grpc.events.listen({
      since: toTimestamp(subMinutes(new Date(), since)),
      until: toTimestamp(new Date()),
    });
    stream.onData((message) => {
      console.log("data");
      requestIdleCallback(() => {
        setEvents((events) => [message.event!].concat(events));
      });
    });
    stream.onEnd(() => console.info("event stream ended"));
    stream.onError((error) => console.error("event stream error: ", error));
    stream.onStatus((status) => console.info("event stream status: ", status));
    return () => stream.cancel();
  }, [since]);

  return (
    <LayoutContent>
      <RadioGroup
        label="Since"
        onChange={(value) => setSince(Number(value.currentTarget.value))}
        selectedValue={since}
      >
        <Radio label="30 Minutes" value={30} />
        <Radio label="1 Hour" value={60} />
        <Radio label="2 Hour" value={120} />
      </RadioGroup>
      <div>
        <Tag>total count = {events.length}</Tag>
      </div>
      <EventList events={events} />
    </LayoutContent>
  );
}
