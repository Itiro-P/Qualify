"use client";

import { Footer, Header } from "@/components";
import { useRouter } from "next/navigation";
import { useState, useEffect } from "react";
import { Analyst } from "@/types/services";
import { analystService } from "@/libs";
import { getSessionUser } from "@/libs/session";
import { Loading } from "@/components/ui";
import { ListReviews } from "@/components/user/listReviews/ListReviews";

export default function Edit() {
  const router = useRouter();
  const [user, setUser] = useState<Analyst | null>(null);

  useEffect(() => {
    async function fetchAnalyst() {
      const session = await getSessionUser();

      if (!session) {
        router.push("/user/register");
        return;
      }
      const user = await analystService.getByUserId(session.id);
      if (user == null) {
        router.push("/user/login");
      } else {
        setUser(user);
      }
    }

    fetchAnalyst();
  }, [router]);

  return (
    <div className="relative flex min-h-screen w-full flex-col overflow-x-hidden">
      <Header />
      <main className="flex-1">
        {user != null ? <ListReviews user={user} /> : <Loading />}
      </main>
      <Footer />
    </div>
  );
}
