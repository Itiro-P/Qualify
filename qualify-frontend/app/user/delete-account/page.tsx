import { DeleteAccount } from "@/components/user";
import { Header, Footer } from "@/components";

export default function DeleteAccountPage() {
  return (
    <div className="relative flex min-h-screen w-full flex-col overflow-x-hidden">
      <Header />
      <main className="flex-1">
        <DeleteAccount />
      </main>
      <Footer />
    </div>
  );
}
