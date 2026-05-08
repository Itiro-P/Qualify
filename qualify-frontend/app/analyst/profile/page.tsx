import { Footer, Header } from "@/components";
import {
  ImageProfile,
  StatisticProfile,
  About,
  Informations,
  ContactButtons,
} from "@/components/analyst/profile";

import { getSessionUser } from "@/libs/session";
import { useState, useEffect } from "react";
import { Analyst } from "@/types/services";
import { analystService } from "@/libs";
import { useRouter } from "next/router";

export default function Profile() {
  const router = useRouter();
  const [analyst, setAnalyst] = useState<Analyst | null>(null);

  useEffect(() => {
    const session = getSessionUser();

    if (!session) {
      router.push("/User/register");
      return;
    }

    analystService
      .getByUserId(session.id)
      .then((resp) => setAnalyst(resp.analyst));
  }, [router]);

  return analyst ? (
    <section id="profile" className="px-6 md:px-20 py-14">
      <Header />

      <div className="flex justify-between ml-3 mt-10">
        <ImageProfile />
        <Informations />
        <ContactButtons />
      </div>

      <div className="flex mt-10">
        <About analyst={analyst} />
        <StatisticProfile />
      </div>

      <Footer />
    </section>
  ) : (
    <p>Usuário não logado ou não é analista</p>
  );
}