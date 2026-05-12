"use client";

import { Footer, Header } from "@/components";
import { useRouter } from "next/navigation";
import {
  ImageProfile,
  StatisticProfile,
  About,
  Informations,
  ContactButtons,
} from "@/components/analyst/profile";
import { useState, useEffect } from "react";
import { Analyst } from "@/types/services";
import { analystService } from "@/libs";
import { getSessionUser } from "@/libs/session";

export default function Profile() {
  const router = useRouter();
  const [analyst, setAnalyst] = useState<Analyst | null>(null);

  async function getAnalyst(id_user: number): Promise<Analyst> {
    return (await analystService.getByUserId(id_user)).analyst;
  }

  useEffect(() => {
    async function fetchAnalyst() {
      const session = getSessionUser();

      if (!session) {
        router.push("/user/register");
        return;
      }
      console.log(session);
      try {
        const analyst = await getAnalyst(session.id);
        setAnalyst(analyst);
      } catch (error) {
        console.error("Erro ao buscar analyst:", error);
      }
    }

    fetchAnalyst();
    console.log(analyst);
    console.log("asdadsadsadasd");
  }, []);

  return (
    <section id="profile" className="px-6 md:px-20 py-14">
      <Header />
      {analyst ? (
        <div>
          <div className="flex justify-between ml-3 mt-10">
            <ImageProfile />
            <Informations />
            <ContactButtons />
          </div>

          <div className="flex mt-10">
            <About analyst={analyst} />
            <StatisticProfile />
          </div>
        </div>
      ) : (
        <p>Usuário não logado ou não é analista</p>
      )}

      <Footer />
    </section>
  );
}
