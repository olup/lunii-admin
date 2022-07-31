import { Box, Button, Center, Icon, Text } from "@chakra-ui/react";
import { BiPlug } from "react-icons/bi";
import { Route } from "wouter";
import { DetailsModal } from "./components/DetailsModal";
import { AboutModal } from "./components/AboutModal";
import { Header } from "./components/Header";
import { NewPackModal } from "./components/NewPackModal";
import { PackList } from "./components/PackList";
import { UpdateModal } from "./components/UpdateModal";
import { useDeviceQuery } from "./hooks/useApi";
import { FiRefreshCcw } from "react-icons/fi";

function App() {
  const {
    data: device,
    isLoading: isLoadingDevice,
    refetch: refetchDevice,
  } = useDeviceQuery();

  if (isLoadingDevice) return null;

  return (
    <Box id="App">
      <Route path="/create-pack">
        <NewPackModal />
      </Route>
      <Route path="/about">
        <AboutModal />
      </Route>
      <Route path="/pack/:id">
        <DetailsModal />
      </Route>
      <Route path="/update">
        <UpdateModal />
      </Route>

      <Box pt="56px" height="100vh">
        <Header />
        {device && (
          <Box p={3}>
            <PackList />
          </Box>
        )}

        {!device && (
          <Center fontSize={20} h="100%" flexDirection="column">
            <Icon as={BiPlug} opacity={0.3}></Icon>
            <Text opacity={0.3} mb={4}>
              No device connected
            </Text>
            <Button
              rightIcon={<FiRefreshCcw />}
              onClick={() => refetchDevice()}
            >
              Refresh
            </Button>
          </Center>
        )}
      </Box>
    </Box>
  );
}

export default App;
