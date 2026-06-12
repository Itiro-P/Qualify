"use client";

import { Suspense } from "react";
import { Footer, Header } from "@/components";
import { SearchAnalyst } from "@/components/customer/searchAnalyst";
import { Loading } from "@/components/ui";

export default function SearchAnalystPage() {
  return (
    <div className="relative flex min-h-screen w-full flex-col overflow-x-hidden">
      <Header />
      <main className="flex-1">
        <Suspense fallback={<Loading />}>
          <SearchAnalyst />
        </Suspense>
      </main>
      <Footer />
    </div>
  );
}
