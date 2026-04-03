import { cookies } from "next/headers";
import { redirect } from "next/navigation";
import { GlobalMap } from "@/components/main/global-map";
import { Hero } from "@/components/main/hero";
import { ProductShowcase } from "@/components/main/product-showcase/index";
import { Stats } from "@/components/main/stats";
import { queryUserInfo } from "@/services/user-api/sdk.gen";
import { NEXT_PUBLIC_API_URL, NEXT_PUBLIC_SITE_URL } from "@/config/constants";

export default async function Home() {
  const Authorization = (await cookies()).get("Authorization")?.value;

  if (Authorization) {
    let user = null;
    try {
      const { data } = await queryUserInfo({
        baseUrl: NEXT_PUBLIC_API_URL || NEXT_PUBLIC_SITE_URL || "",
        headers: {
          Authorization,
        },
      });
      user = data;
    } catch (error) {
      console.log("Token validation failed:", error);
    }

    if (user) {
      redirect("/dashboard");
    }
  }

  return (
    <main className="container space-y-16">
      <Hero />
      <Stats />
      <ProductShowcase />
      <GlobalMap />
    </main>
  );
}
