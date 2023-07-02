# v0.3.6 (pre
* feature
    - The way of previewing report has been changed from clicking tags to displaying directly at the bottom, with tab switching for charts (more intuitive) 
    - Replaced the implementation library of codemirror to provide a better code writing experience
    - Added share feature,  by selecting bot in bots panel and clicking share can copy the bot's address to clipboard for others to access directly the bot's editing view
    - Added automatic refresh for running (default 10s)  
    - Changed the storage implementation of batch, now it will be stored in the db so that it can continue executing after an abnormal interruption
* fix   
    - It will directly panic if the database cannot connect (encounter errors should terminate immediately)      
    - Replaced the clipboard implementation of share button to an earlier api (can adapt to more browsers)   
    - Fixed the problem that report was not sorted by time   
    - Fixed the issue of wrong click event in bots

# v0.3.5 (pre
Major adjustments are nearing completion. Version 0.3 will only fix bugs next.
* Feature
    - Added a queue delay configuration to control the scheduling frequency of the robot
    - Optimized the CSS implementation of nodes in the sideplane
    - Added HTTP query params as input
* Fix
    - The code input box has disordered input logic after switching input methods
    - Clicking inputnumber in bots will lose focus
    - When zooming and resizing the window, the editor window is not enlarged or reduced proportionally
    
# v0.3.1 (pre
* Feature
    - Added a button to erase the behavior tree
    - The drawing of the graph now depends entirely on the data in model/tree
* Fix
    - Clicking too fast caused the current node to draw incorrectly
    - Some jumps in the debugging window are fixed
    - When zooming and resizing the window, the editor window is not enlarged or reduced proportionally
# v0.3.0 (pre

* Feature
    - Rewrote the entire editor using the umi framework and ts (type safe, supports dark mode switching
    - Replaced components with functions, wrote code using hooks + redux (stateless mode (optimized loading time and drawing efficiency
    - Added a new bot loading method that can load a bot by accessing the URL (easier to spread
    - Introduced the SQLite in-memory database (easy to try, can be deployed locally quickly
* Fix
    - Fixed the loss of tail node information
    - Fixed the problem that batches could not exit accurately during pressure testing

# v0.2.5 
* Feature
    - Rewrote the sidebar to provide a better filtering method
    - Prefab is displayed separately as a page and provides search and editing functions
    - Optimized connection points (shrink when there is no mouse movein)
    - Added time sorting to the report page

## v0.2.1
* Feature
   - Added new parallel nodes
   - Deleted the original assert node type
   - Added a runtime err column to output runtime error information
   - Refactored the runtime logic of bots
   - Introduced the logic to display thread information in the response column (parallel nodes will create new threads)
   - Added a small animation when running to the node (optimization prompt

## v0.1.17 
* Feature
* Fix
    - Fixed logic errors in asynchronous loading of behavior trees

## v0.1.16
* Feature
    - Provides fmt function for lua code
    - Leave enough space for prefab and move change to overlap with meta window
    - Add step shortcut [F10]
* Fix
    - Error deleting configuration file failed
    - Root node position correction
    - Fix health check not performed after initializing server address

## v0.1.15
* Feature
   - Removed the debug button in the edit interface and automatically created when clicking step 
   - Added a delay for step on the last node (to prevent continuous clicking 
   - Added a reset button to prevent the user from not wanting to execute down
   - Removed the utils in the script module (changed to an independent uuid and random interface displayed at the first level directory for easy reference 
   - Added node "prefab" function (now users can define and reuse their own script nodes in the config panel 
   - Added connection status prompt
* Fix
    - step api does not return correct error information