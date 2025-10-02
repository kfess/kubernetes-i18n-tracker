import { createBrowserRouter, RouterProvider } from 'react-router-dom';
import { NavigationBar } from '@/features/NavigationBar';
import { DetailPage } from './pages/DetailPage';
import { DevPage } from './pages/DevPage';
import { HomePage } from './pages/Home.page';

const basename = import.meta.env.MODE === 'production' ? '/kubernetes-i18n-tracker' : '/';

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
      path: '/detail/:id',
      element: (
        <NavigationBar>
          <DetailPage />
        </NavigationBar>
      ),
    },
    {
      path: '/dev',
      element: (
        <NavigationBar>
          <DevPage />
        </NavigationBar>
      ),
    },
  ],
  { basename }
);

export function Router() {
  return <RouterProvider router={router} />;
}
