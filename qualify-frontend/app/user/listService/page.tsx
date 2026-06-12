"use client";

import { Footer, Header } from "@/components";
import { useRouter } from "next/navigation";
import { useState, useEffect } from "react";
import { getSessionUser, User } from "@/libs/session";
import { Loading } from "@/components/ui";
import { ListServices } from "@/components/user/listServices/ListServices";

export default function Edit() {
  const router = useRouter();
  const [user, setUser] = useState<User | null>(null);

  useEffect(() => {
    async function fetchUser() {
      const session = await getSessionUser();

      if (!session) {
        console.log("No session found, redirecting to register");
        router.push("/user/register");
        return;
      }
      setUser(session);
    }

    fetchUser();
  }, [router]);

  return (
    <div className="relative flex min-h-screen w-full flex-col overflow-x-hidden">
      <Header />
      <main className="flex-1">
        {user != null ? <ListServices user={user} /> : <Loading />}
      </main>
      <Footer />
    </div>
  );
}
