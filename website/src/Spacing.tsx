import styled from 'styled-components';

interface Props {
  spacing: number;
  horizontal?: boolean;
  justify?: 'left' | 'right' | 'center';
  padding?: number | string;
}

export const Spacing = styled.div<Props>`
  display: grid;
  grid-gap: ${(props) => props.spacing}rem;
  grid-template: auto/ auto;
  grid-auto-flow: ${(props) => (props.horizontal ? 'column' : 'row')};
  justify-content: ${(props) => (props.horizontal ? 'start' : 'normal')};
  justify-items: ${(props) => props.justify ?? 'normal'};
  align-items: flex-start;
  padding: ${(props) => props.padding ?? 0};
`;
