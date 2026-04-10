import {
  DeviceTesting,
  Footer,
  Header,
  HeroSection,
  RetestGuarantee,
  TestingTiers,
} from "@/components/landing";

import {
  Perfil,
} from "@/components/perfil_analista";

export default function Home() {
  return (
    <div className="relative flex min-h-screen w-full flex-col overflow-x-hidden">
      <Header />
      <main className="flex-1">
        <HeroSection />
        <TestingTiers />
        <DeviceTesting />
        <RetestGuarantee />
        <Perfil />
      </main>
      <Footer />
    </div>
  );
}
