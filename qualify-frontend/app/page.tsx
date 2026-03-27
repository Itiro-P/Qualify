import {
  Header,
  HeroSection,
  TestingTiers,
  DeviceTesting,
  RetestGuarantee,
  Footer,
} from "@/components/landing";

export default function Home() {
  return (
    <div className="relative flex min-h-screen w-full flex-col overflow-x-hidden">
      <Header />
      <main className="flex-1">
        <HeroSection />
        <TestingTiers />
        <DeviceTesting />
        <RetestGuarantee />
      </main>
      <Footer />
    </div>
  );
}
