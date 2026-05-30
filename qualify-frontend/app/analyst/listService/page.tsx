"use client";

import { Footer, Header } from "@/components";
import { useRouter } from "next/navigation";
import { useState, useEffect } from "react";
import { Analyst } from "@/types/services";
import { analystService } from "@/libs";
import { getSessionUser } from "@/libs/session";
import { Loading } from "@/components/ui";
import { ListServices } from "@/components/analyst/listServices/ListServices";

export default function Edit() {
  const router = useRouter();
  const [analyst, setAnalyst] = useState<Analyst | null>(null);

  useEffect(() => {
    async function fetchAnalyst() {
      const session = await getSessionUser();

      if (!session) {
        router.push("/user/register");
        return;
      }
      const analyst = await analystService.getByUserId(session.id);
      if (analyst == null) {
        router.push("/user/login");
      } else {
        setAnalyst(analyst);
      }
    }

    fetchAnalyst();
  }, [router]);

  return (
    <div>
      <Header />
      {analyst != null ? <ListServices analyst={analyst} /> : <Loading />}
      <Footer />
    </div>
  );
}
