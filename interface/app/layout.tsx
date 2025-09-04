import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "WebVM",
  description: "Launch a temporary virtual machine instantly in your browser - no setup required!",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body style={{ minHeight: "100vh" }}>
        {children}
      </body>
    </html>
  );
}
