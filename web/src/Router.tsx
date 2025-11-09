import { createBrowserRouter, RouterProvider } from 'react-router-dom';
import { NavigationBar } from '@/features/NavigationBar';
import { basename } from './const';
import { DetailPage } from './pages/DetailPage';
import { HomePage } from './pages/Home.page';
import { InformationPage } from './pages/InformationPage';

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
