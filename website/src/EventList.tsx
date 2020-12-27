import { Button, Classes, Dialog, H5, Tag } from "@blueprintjs/core";
import { formatDistanceToNow } from "date-fns";
import React from "react";
import styled from "styled-components";
import { toDate } from "./grpc";
import { present } from "./Present";
import { Event } from "./sdk/events_pb";
import { Spacing } from "./Spacing";

interface Props {
  events: Event.AsObject[];
}

export function EventList(props: Props) {
  return (
    <Spacing spacing={0}>
      {props.events.map((event) => (
        <EventRow key={event.id} event={event} />
      ))}
    </Spacing>
  );
}

const EventRow = React.memo((props: { event: Event.AsObject }) => {
  return (
    <Row onClick={() => presentEventDialog(props.event)}>
      <Priority event={props.event} />
      <div>
        <H5>{props.event.rule}</H5>
        <EventInfo event={props.event} />
      </div>
      <div>{formatDistanceToNow(toDate(props.event.time!))} ago</div>
    </Row>
  );
});

const Row = styled.div`
  display: grid;
  grid-template-columns: min-content auto max-content;
  grid-template-rows: 1fr;
  gap: 2rem;
  padding: 1rem 2rem;
  border-bottom: 2px solid rgb(211, 211, 211);
  &:hover {
    background: rgb(189, 225, 244);
    cursor: pointer;
  }
`;

const Priority = styled.div<{ event: Event.AsObject }>`
  align-self: center;
  justify-self: center;
  width: 13px;
  height: 13px;
  border-radius: 50%;
  background: ${(props) => {
    switch (props.event.priority) {
      case "Error":
        return "red";
      case "Warning":
        return "orange";
      case "Notice":
        return "blue";
      case "Debug":
      // fallthrough
      default:
        return "grey";
    }
  }};
`;

function EventInfo(props: { event: Event.AsObject }) {
  const fields = React.useMemo(
    () => JSON.parse(props.event.raw)["output_fields"],
    [props.event.raw]
  );

  return (
    <Spacing spacing={0.3} justify="left">
      <FieldTag label="namespace" value={fields["k8s.ns.name"]} />
      <FieldTag label="pod" value={fields["k8s.pod.name"]} />
      <FieldTag label="process" value={fields["proc.name"]} />
      <FieldTag label="parent process" value={fields["proc.pname"]} />
      <FieldTag label="file descriptor" value={fields["fd.name"]} />
    </Spacing>
  );
}

function FieldTag(props: { label: string; value: any }) {
  if (props.value === undefined || props.value === null) {
    return null;
  }
  return (
    <Tag>
      {props.label} = {props.value}
    </Tag>
  );
}

const presentEventDialog = (event: Event.AsObject) => {
  present(({ close }) => (
    <Dialog title="Falco Event" icon="info-sign" isOpen={true} onClose={close}>
      <div className={Classes.DIALOG_BODY}>
        <pre className={Classes.CODE_BLOCK} style={{ overflowX: "scroll" }}>
          <code>{JSON.stringify(JSON.parse(event.raw), null, "  ")}</code>
        </pre>
      </div>
      <div className={Classes.DIALOG_FOOTER}>
        <div className={Classes.DIALOG_FOOTER_ACTIONS}>
          <Button onClick={close}>Close</Button>
        </div>
      </div>
    </Dialog>
  ));
};
