import { createBrowserRouter, RouterProvider } from 'react-router-dom';
import { NavigationBar } from '@/features/NavigationBar';
import { basename } from './const';
import { DetailPage } from './pages/DetailPage';
import { DevPage } from './pages/DevPage';
import { HomePage } from './pages/Home.page';

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
