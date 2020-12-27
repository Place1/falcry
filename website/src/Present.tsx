import React from "react";
import { unmountComponentAtNode, render } from "react-dom";

interface ContentProps {
  close: () => void;
}

export function present(Content: React.FC<ContentProps>) {
  requestAnimationFrame(() => {
    const root = document.createElement("div");
    document.body.appendChild(root);
    const close = () => unmountComponentAtNode(root);
    render(<Content close={close} />, root);
  });
}
