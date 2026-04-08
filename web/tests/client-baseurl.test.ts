import { describe, expect, test } from "bun:test";

import "../apps/admin/utils/setup-clients";
import "../apps/user/utils/setup-clients";
import { client as adminAdminClient } from "../apps/admin/services/admin-api/client.gen";
import { client as adminCommonClient } from "../apps/admin/services/common-api/client.gen";
import { client as adminUserClient } from "../apps/admin/services/user-api/client.gen";
import { client as userCommonClient } from "../apps/user/services/common-api/client.gen";
import { client as userUserClient } from "../apps/user/services/user-api/client.gen";

describe("client base urls", () => {
  test("user app keeps common-api on /api/v1/common and user-api at site root", () => {
    expect(userCommonClient.getConfig().baseUrl).toBe("/api/v1/common");
    expect(userUserClient.getConfig().baseUrl).toBe("");
  });

  test("admin app keeps admin/common prefixes and user-api at site root", () => {
    expect(adminCommonClient.getConfig().baseUrl).toBe("/api/v1/common");
    expect(adminAdminClient.getConfig().baseUrl).toBe("/api/v1/admin");
    expect(adminUserClient.getConfig().baseUrl).toBe("");
  });
});
