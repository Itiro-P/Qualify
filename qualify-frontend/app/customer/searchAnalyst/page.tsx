"use client";

import { Footer, Header } from "@/components";
import { SearchAnalyst } from "@/components/customer/searchAnalyst";

export default function Profile() {
  return (
    <section id="profile" className="px-6 md:px-20 py-14">
      <Header />
      <SearchAnalyst />
      <Footer />
    </section>
  );
}
