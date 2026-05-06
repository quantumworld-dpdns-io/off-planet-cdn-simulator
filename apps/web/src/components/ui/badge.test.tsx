import { render, screen } from "@testing-library/react";
import { Badge } from "./badge";

describe("Badge", () => {
  it("renders children", () => {
    render(<Badge variant="p0">P0</Badge>);
    expect(screen.getByText("P0")).toBeInTheDocument();
  });

  it("applies p0 red styles", () => {
    const { container } = render(<Badge variant="p0">Emergency</Badge>);
    expect(container.firstChild).toHaveClass("bg-red-100");
  });

  it("applies online green styles", () => {
    const { container } = render(<Badge variant="online">Online</Badge>);
    expect(container.firstChild).toHaveClass("bg-green-100");
  });

  it("applies p4 grey styles", () => {
    const { container } = render(<Badge variant="p4">Entertainment</Badge>);
    expect(container.firstChild).toHaveClass("bg-gray-100");
  });
});
