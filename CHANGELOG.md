# v0.4.3
* Features
    - Renamed the meta panel to blackboard, aligning with the semantics of behavior trees #19
    - Fixed the length of delayed label display to avoid affecting the layout display of other controls under different delays #15
    - Changed the original memory type from SQLite to file-based SQLite (preventing loss of robot files upon restart) #17
* Fixes
    - Fixed the duplicate construction error in heartbeat check delay (causing excessive refresh) #18

# v.0.4.2
* Features
    - The message module is provided, and users can now process stream byte data by themselves at the script layer (unpacking, packaging
    - The concept of report has been modified, and information such as request time-consuming is no longer provided (it is more reasonable to leave it to the background for statistics). Report currently only provides statistics on the number dimensions such as req, res, ntf and so on.

# v0.4.1
* Features
    - Added websocket module
* Fixes
    - Fixed missing banner print

# v0.4.0
* Features
    - Added cluster deployment mode

# v0.3.6 (pre
* Features
    - The way of previewing report has been changed from clicking tags to displaying directly at the bottom, with tab switching for charts (more intuitive) 
    - Replaced the implementation library of codemirror to provide a better code writing experience
    - Added share Features,  by selecting bot in bots panel and clicking share can copy the bot's address to clipboard for others to access directly the bot's editing view
    - Added automatic refresh for running (default 10s)  
    - Changed the storage implementation of batch, now it will be stored in the db so that it can continue executing after an abnormal interruption
* Fixes   
    - It will directly panic if the database cannot connect (encounter errors should terminate immediately)      
    - Replaced the clipboard implementation of share button to an earlier api (can adapt to more browsers)   
    - Fixed the problem that report was not sorted by time   
    - Fixed the issue of wrong click event in bots

# v0.3.5 (pre
Major adjustments are nearing completion. Version 0.3 will only Fixes bugs next.
* Features
    - Added a queue delay configuration to control the scheduling frequency of the robot
    - Optimized the CSS implementation of nodes in the sideplane
    - Added HTTP query params as input
* Fixes
    - The code input box has disordered input logic after switching input methods
    - Clicking inputnumber in bots will lose focus
    - When zooming and resizing the window, the editor window is not enlarged or reduced proportionally
    
# v0.3.1 (pre
* Features
    - Added a button to erase the behavior tree
    - The drawing of the graph now depends entirely on the data in model/tree
* Fixes
    - Clicking too fast caused the current node to draw incorrectly
    - Some jumps in the debugging window are fixed
    - When zooming and resizing the window, the editor window is not enlarged or reduced proportionally
# v0.3.0 (pre

* Features
    - Rewrote the entire editor using the umi framework and ts (type safe, supports dark mode switching
    - Replaced components with functions, wrote code using hooks + redux (stateless mode (optimized loading time and drawing efficiency
    - Added a new bot loading method that can load a bot by accessing the URL (easier to spread
    - Introduced the SQLite in-memory database (easy to try, can be deployed locally quickly
* Fixes
    - Fixed the loss of tail node information
    - Fixed the problem that batches could not exit accurately during pressure testing

# v0.2.5 
* Features
    - Rewrote the sidebar to provide a better filtering method
    - Prefab is displayed separately as a page and provides search and editing functions
    - Optimized connection points (shrink when there is no mouse movein)
    - Added time sorting to the report page

## v0.2.1
* Features
   - Added new parallel nodes
   - Deleted the original assert node type
   - Added a runtime err column to output runtime error information
   - Refactored the runtime logic of bots
   - Introduced the logic to display thread information in the response column (parallel nodes will create new threads)
   - Added a small animation when running to the node (optimization prompt

## v0.1.17 
* Features
* Fixes
    - Fixed logic errors in asynchronous loading of behavior trees

## v0.1.16
* Features
    - Provides fmt function for lua code
    - Leave enough space for prefab and move change to overlap with meta window
    - Add step shortcut [F10]
* Fixes
    - Error deleting configuration file failed
    - Root node position correction
    - Fixes health check not performed after initializing server address

## v0.1.15
* Features
   - Removed the debug button in the edit interface and automatically created when clicking step 
   - Added a delay for step on the last node (to prevent continuous clicking 
   - Added a reset button to prevent the user from not wanting to execute down
   - Removed the utils in the script module (changed to an independent uuid and random interface displayed at the first level directory for easy reference 
   - Added node "prefab" function (now users can define and reuse their own script nodes in the config panel 
   - Added connection status prompt
* Fixes
    - step api does not return correct error information