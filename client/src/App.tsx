import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import Layout from "./components/code/Layout";
import Dashboard from "./pages/Dashboard";
import Accounts from "./pages/Accounts";
import Employees from "./pages/Employees";
import Payroll from "./pages/Payroll";
import Withdrawals from "./pages/Withdrawals";
import History from "./pages/History";

const queryClient = new QueryClient();

const App = () => {
  return (
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<Layout />}>
            <Route index element={<Dashboard />} />
            <Route path="/accounts" element={<Accounts />} />
            <Route path="/employees" element={<Employees />} />
            <Route path="/payroll" element={<Payroll />} />
            <Route path="/withdrawals" element={<Withdrawals />} />
            <Route path="/history" element={<History />} />
          </Route>
        </Routes>
      </BrowserRouter>
    </QueryClientProvider>
  );
};

export default App;
