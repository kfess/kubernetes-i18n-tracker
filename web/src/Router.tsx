import { createBrowserRouter, RouterProvider } from 'react-router-dom';
import { NavigationBar } from '@/features/NavigationBar';
import { basename } from './const';
import { DetailPage } from './pages/DetailPage';
import { HomePage } from './pages/Home.page';
import { InformationPage } from './pages/InformationPage';
import { SummaryPage } from './pages/SummaryPage';

const router = createBrowserRouter(
  [
    {
      path: '/',
      element: (
        <NavigationBar>
          <HomePage />
        </NavigationBar>
      ),
    },
    {
      path: '/detail',
      element: (
        <NavigationBar>
          <DetailPage />
        </NavigationBar>
      ),
    },
    {
      path: '/summary',
      element: (
        <NavigationBar>
          <SummaryPage />
        </NavigationBar>
      ),
    },
    // {
    //   path: '/information',
    //   element: (
    //     <NavigationBar>
    //       <InformationPage />
    //     </NavigationBar>
    //   ),
    // },
  ],
  { basename }
);

export function Router() {
  return <RouterProvider router={router} />;
}
