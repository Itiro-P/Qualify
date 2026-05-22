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
import { Loading } from "@/components/ui";

export default function Profile() {
  const router = useRouter();
  const [analyst, setAnalyst] = useState<Analyst | null>(null);

  useEffect(() => {
    async function fetchAnalyst() {
      const session = getSessionUser();

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
  }, []);

  return (
    <section id="profile" className="px-6 md:px-20 py-14">
      <Header />
      {analyst != null ? (
        <div>
          <div className="flex justify-between ml-3 mt-10">
            <ImageProfile />
            <Informations
              name={analyst.name}
              city={analyst.city}
              state={analyst.country_state}
              country={analyst.country_name}
              rating={analyst.mean_rating}
              reviews={analyst.total_reviews}
              date={analyst.time_created}
            />
            <ContactButtons />
          </div>

          <div className="flex mt-10">
            <About analyst={analyst} />
            <StatisticProfile
              analyst_id={analyst.id}
              hourly_rate={analyst.hourly_rate}
            />
          </div>
        </div>
      ) : (
        <Loading />
      )}

      <Footer />
    </section>
  );
}
