import { describe, expect, test } from "bun:test";
import { renderToStaticMarkup } from "react-dom/server";
import { MemoryRouter } from "react-router-dom";

import { SidebarProvider } from "../../../packages/ui/src/components/sidebar";
import { CollapsedStandaloneNavLink } from "../components/collapsed-standalone-nav-link";

describe("admin collapsed navigation", () => {
  test("renders standalone collapsed nav items as a real anchor instead of button nesting", () => {
    const html = renderToStaticMarkup(
      <MemoryRouter>
        <SidebarProvider defaultOpen={false}>
          <CollapsedStandaloneNavLink
            href="/dashboard/workplace"
            label="仪表盘"
            isActive={false}
            icon={<span data-testid="icon" />}
          />
        </SidebarProvider>
      </MemoryRouter>,
    );

    expect(html).toContain('<a data-sidebar="menu-button"');
    expect(html).toContain('href="/admin/dashboard/workplace"');
    expect(html).not.toContain("<button");
  });
});
