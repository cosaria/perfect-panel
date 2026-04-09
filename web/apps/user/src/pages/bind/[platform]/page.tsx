import BindContent from "./bind-content";

export function generateStaticParams() {
  return [
    { platform: "telegram" },
    { platform: "apple" },
    { platform: "facebook" },
    { platform: "google" },
    { platform: "github" },
  ];
}

export default function Page() {
  return <BindContent />;
}
