import { Flex, IconButton, Box } from "@chakra-ui/react";
import { useQuery } from "@tanstack/react-query";

import { DragDropContext, Droppable, Draggable } from "@hello-pangea/dnd";
import { FiMoreVertical } from "react-icons/fi";
import { MdDragIndicator } from "react-icons/md";
import { Link } from "wouter";
import { ListPacks } from "../../wailsjs/go/main/App";
import { useDeviceQuery } from "../hooks/useApi";
import { useChangePackOrder } from "../hooks/useChangePackOrder";
import { PackTag } from "./PackTag";

export const PackList = () => {
  const { data: device } = useDeviceQuery();
  const { mutate: handleChangePackOrder } = useChangePackOrder();

  const { data: packs } = useQuery(["packs"], ListPacks, {
    enabled: !!device,
  });

  return (
    <Box>
      <Flex py={3}>
        <Box width={50}></Box>
        <Box flex={1}>Title</Box>
        <Box width={100}></Box>
        <Box width={50}></Box>
      </Flex>
      <DragDropContext
        onDragEnd={(a) => {
          if (!a.destination) return;
          handleChangePackOrder({
            id: a.draggableId,
            position: a.destination?.index,
          });
        }}
      >
        <Droppable droppableId="droppable">
          {(provided) => (
            <Box {...provided.droppableProps} ref={provided.innerRef}>
              {packs?.map((p, i) => (
                <Draggable
                  key={p.uuid as any}
                  draggableId={p.uuid as any}
                  index={i}
                >
                  {(provided) => (
                    <Flex
                      ref={provided.innerRef}
                      {...provided.draggableProps}
                      backgroundColor="white"
                      borderBottom="1px solid #eee"
                      borderTop="1px solid #eee"
                      mt="-1px"
                      py={2}
                      alignItems="center"
                    >
                      <Box width={50}>
                        <IconButton
                          size="xs"
                          aria-label="up"
                          icon={<MdDragIndicator />}
                          variant="ghost"
                          mr={1}
                          {...provided.dragHandleProps}
                        />
                      </Box>

                      <Box
                        flex={1}
                        fontWeight={p.title && "bold"}
                        opacity={p.title ? 1 : 0.5}
                      >
                        {p.title || p.uuid}
                      </Box>
                      <Box width={150}>
                        <PackTag metadata={p} />
                      </Box>
                      <Box width={50}>
                        <Link to={`/pack/${p.uuid}`}>
                          <IconButton
                            variant="ghost"
                            aria-label="Details"
                            icon={<FiMoreVertical />}
                          />
                        </Link>
                      </Box>
                    </Flex>
                  )}
                </Draggable>
              ))}
            </Box>
          )}
        </Droppable>
      </DragDropContext>
    </Box>
  );
};
